package main

import (
	//"fmt"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/RomanLevBy/UrlShortener/internal/config"
	"github.com/RomanLevBy/UrlShortener/internal/http-server/handlers/redirect"
	"github.com/RomanLevBy/UrlShortener/internal/http-server/handlers/url/delete"
	"github.com/RomanLevBy/UrlShortener/internal/http-server/handlers/url/save"
	mwLogger "github.com/RomanLevBy/UrlShortener/internal/http-server/middleware/logger"
	"github.com/RomanLevBy/UrlShortener/internal/lib/logger/sl"
	"github.com/RomanLevBy/UrlShortener/internal/storage/postgres"
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
	conf := config.MustLoad()

	log := setupLogger(conf.Env)
	log.Info("Config init.", slog.String("env", conf.Env))
	log.Debug("Debug messages are enable")

	storage, err := postgres.New(conf)
	if err != nil {
		log.Error("Failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	//router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/api/url", func(r chi.Router) {
		//todo create new auth approach
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			conf.HTTPServer.User: conf.HTTPServer.Password,
		}))

		r.Post("/", save.New(log, storage))
		r.Delete("/{alias}", delete.New(log, storage))
	})

	router.Get("/{alias}", redirect.New(log, storage))

	log.Info("starting server", slog.String("address", conf.Address))

	srv := &http.Server{
		Addr:         conf.Address,
		Handler:      router,
		ReadTimeout:  conf.HTTPServer.Timeout,
		WriteTimeout: conf.HTTPServer.Timeout,
		IdleTimeout:  conf.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server", err)
	}

	log.Info("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
