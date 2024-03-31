package redirect

import (
	"log/slog"

	"errors"
	"net/http"

	resp "github.com/RomanLevBy/UrlShortener/internal/lib/api/response"
	"github.com/RomanLevBy/UrlShortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type URLProvider interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlProvider URLProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.redirect.New"

		log.With(
			slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("Alias is empty", slog.String("alias", alias))

			render.JSON(w, r, resp.Error("alias is empty"))

			return
		}

		url, err := urlProvider.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias", alias)

			render.JSON(w, r, resp.Error("not found"))

			return
		}

		if err != nil {
			log.Info("Fail to get url", slog.String("alias", alias))

			render.JSON(w, r, resp.Error("fail to get url"))

			return
		}

		log.Info("Redirect to url", slog.String("url", url))

		http.Redirect(w, r, url, http.StatusSeeOther)
	}
}
