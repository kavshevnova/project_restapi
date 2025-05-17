package save

import (
	"errors"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	resp "github.com/kavshevova/project_restapi/internal/lib/api/response"
	"github.com/kavshevova/project_restapi/internal/lib/logger/sl"
	"github.com/kavshevova/project_restapi/internal/lib/random"
	"github.com/kavshevova/project_restapi/internal/storage"
	"log/slog"
	"net/http"
)

//запросы, которые будут поступать в виде джсона
type Request struct {
	URL string `json:"url" validate:"required,url"`
	Alias string `json:"alias"`
}

type Response struct {
	resp.Response
	Alias  string `json:"alias,omitempty"` //элиас только что сохраненного урла
}

//TODO: перенести алиас в конфиг или в базу данных
const aliasLength = 6

type URLSaver interface {
	SaveURL(urlTOsave string, alias string) (int64, error)
}

//функция-конструктор для хендлера, то есть во время подключения этого хендлера к роутеру мы будем вызываать функцию new, которая возвращает хендлер и здесь мы можем передать доп параметры которые будут установлены в каждом обработчике
func New (log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("requestID", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req) //декодируем запрос в объект реквест
		if err != nil {

			//пишем ошбку в лог
			log.Error("failed to decode request body", sl.Err(err))

			//возвращаем джсон с ответом нашему клиенту
			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("reques body decoded", slog.Any("request", req))

		//создаем новый объект валидатора и валидируем структуру req
		if err := validator.New().Struct(req); err != nil {
			//если валидатор находит ошибку он возвращает ошибку вот такого типа
			validateErr := err.(validator.ValidationErrors)
			//эту ошибку мы в чистом виде залогируем без изменений
			log.Error("invalid request", sl.Err(err))
			//сформируем готовый запрос в который уже вписано человекочитаемые ошибки
			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}
		alias := req.Alias
		if alias == "" {
			//если  алиас пустой то мы генерируем его из случайных символов
			alias = random.NewRandomString(aliasLength)
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url", req.URL))
			render.JSON(w, r, resp.Error("url already exists"))
			return
		}
		if err != nil {
			log.Error("failed to add url", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to add url"))
			return
		}
		log.Info("url added", slog.Int64("id", id))
		ResponseOk(w, r, alias)
	})
}

func ResponseOk(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{
		Response: resp.OK(),
			Alias:    alias,
		})
	}
