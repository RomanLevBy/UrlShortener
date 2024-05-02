package delete

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"

	resp "github.com/RomanLevBy/UrlShortener/internal/lib/api/response"
	"github.com/RomanLevBy/UrlShortener/internal/lib/logger/sl"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.1 --name=URLRemover
type URLRemover interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlRemover URLRemover) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.url.delete.New"

		log = log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty", slog.String("alias", alias))

			render.JSON(w, r, resp.Error("Alias is empty"))

			return
		}

		err := urlRemover.DeleteURL(alias)
		if err != nil {
			log.Info("fail to delete url", slog.String("alias", alias), sl.Err(err))

			render.JSON(w, r, resp.Error("Fail to delete url"))

			return
		}

		log.Info("url deleted", slog.String("alias", alias))

		render.JSON(w, r, resp.OK())
	}
}
