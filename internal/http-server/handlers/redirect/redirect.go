package redirect

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	resp "github.com/kavshevova/project_restapi/internal/lib/api/response"
	"github.com/kavshevova/project_restapi/internal/lib/logger/sl"
	"github.com/kavshevova/project_restapi/internal/storage"
	"log/slog"
	"net/http"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}
//go:generate mockery --name URLGetter --with-expecter
func New(log *slog.Logger, getter URLGetter ) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	const op = "handlers.redirect.New"

	log = log.With(
		slog.String("op", op),
		slog.String("requestID", middleware.GetReqID(r.Context())),
		)

	alias := chi.URLParam(r, "alias")
	if alias == "" {
		log.Info("alias is empty")
		render.JSON(w, r, resp.Error("invalid request"))
		return
	}

	resURL, err := getter.GetURL(alias)
	if errors.Is(err, storage.ErrURLNotFound) {
		log.Info("url not found")
		render.JSON(w, r, resp.Error("not found"))
		return
	}
	if err != nil {
		log.Error("failed to get url", sl.Err(err))
		render.JSON(w, r, resp.Error("internal server error"))
		return
	}
	log.Info("get url", slog.String("url", resURL))

	http.Redirect(w, r, resURL, http.StatusFound)
	})
}

