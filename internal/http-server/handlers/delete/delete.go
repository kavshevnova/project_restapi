package delete

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	resp "github.com/kavshevova/project_restapi/internal/lib/api/response"
	"github.com/kavshevova/project_restapi/internal/lib/logger/sl"
	"log/slog"
	"net/http"
)

type URLDeleter interface {
	DeleteURL(alias string) error
}
//go:generate mockery --name URLDeleter --with-expecter
func New(log *slog.Logger, deleter URLDeleter ) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.delete.New"

		log = log.With(
			slog.String("op", op),
			slog.String("requestID", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")
			w.WriteHeader(http.StatusBadRequest) // 400
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		result := deleter.DeleteURL(alias)
		if result != nil {
			log.Error("failed to delete url", sl.Err(result))
			w.WriteHeader(http.StatusInternalServerError) // 500
			render.JSON(w, r, resp.Error("internal server error"))
			return
		}

		log.Info("url deleted", slog.String("alias", alias))
		render.JSON(w,r, resp.OK())
	})
	}