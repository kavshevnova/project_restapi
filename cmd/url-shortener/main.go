package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/kavshevova/project_restapi/internal/config"
	delete2 "github.com/kavshevova/project_restapi/internal/http-server/handlers/delete"
	"github.com/kavshevova/project_restapi/internal/http-server/handlers/redirect"
	"github.com/kavshevova/project_restapi/internal/http-server/handlers/url/save"
	"github.com/kavshevova/project_restapi/internal/http-server/middleware/logger"
	"github.com/kavshevova/project_restapi/internal/lib/logger/handlers/slogpretty"
	"github.com/kavshevova/project_restapi/internal/lib/logger/sl"
	"github.com/kavshevova/project_restapi/internal/storage/sqlite"
	"log"
	"log/slog"
	"net/http"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	log.Info("Starting server", slog.String("env", cfg.Env))
	log.Debug("debug logging enabled")
	log.Error("ошибка")
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1) //завершаем программу с кодом 1
	}
	_ = storage

	//инициализируем роутер через пакет чи
	router := chi.NewRouter()
	//подключаем к роутеру мидлвейр, мидлвейр это хендлеры в цепочке которые обрабатывают не основной запрос, например хендлер проверяющий авторизацию в цепочке для основного хендлера создание или удаление урла
	router.Use(middleware.RequestID) //суть этого мидлвеера что он добавляет к каждому поступающему запросу уникальный реквестайди для того чтобы если что-то пошло не так в одном запросе можно было его найти по айди и разобрать
	router.Use(middleware.Logger)    //посмотреть айпи пользователя который к нам постучался
	router.Use(logger.New(log))      //логирует все входящие запросы, будет добавлена строчка лог, которая говорит о том что я получил запрос я его обработал и на обработку ушло столько то времени
	router.Use(middleware.Recoverer) //если случается паника внутри хендлера, из-за одного запроса не должно падать все приложение целиком поэтому мы восстанавливаем эту панику
	router.Use(middleware.URLFormat) //чтобы можно было писать красивые урлы при подключении их к обработчику к нашему роутеру

	//делаем роутер внутри роутера
	router.Route("/url", func(r chi.Router) {
		//поддключаем авторизацию
		//basicauth максимально простая авторизация из пакета чи которая предполагает отправку логина и пароля в заголовке
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))
		r.Post("/", save.New(log, storage)) //r.Post так как внутренний роутер r и убираем приставку /url  так как она уже есть у всей группы
		r.Delete("/{alias}", delete2.New(log, storage))
	})

	router.Get("/{alias}", redirect.New(log, storage))

	log.Info("Starting server", slog.String("env", cfg.Address))
	//создаем сам сервер через http библиотеку благодаря совместимости chi с этой библиотекой
	srv := &http.Server{
		Addr:         cfg.Address,            //наш адрес из конфига
		Handler:      router,                 //вся группа хендлеров роутера заключена в роутер который тоже является хендлером
		ReadTimeout:  cfg.HTTPServer.Timeout, //наш таймаут на обработку запросов (время на чтение запроса)
		WriteTimeout: cfg.HTTPServer.Timeout, //наш таймаут на обработку запросов (время на ответ клиенту)
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()

	log.Error("server stopped")

	//TODO: run server
}

func setupLogger(env string) *slog.Logger {
	//почему логгер должен зависеть от параметра енв:
	//локально мы хотим видеть текстовые логи
	//в окружении дев или прод (на сервере) мы хотим видеть json. Причем на дев - логи уровня дебаг, а на проде не ниже уровня инфо.
	var log *slog.Logger
	//объявили логгер и в зависимости от переменной енв мы будем его создавать
	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = setupPrettySlog()
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}
	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}