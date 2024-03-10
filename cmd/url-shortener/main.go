package main

import (
	//"fmt"
	"github.com/RomanLevBy/UrlShortener/internal/config"
	"github.com/RomanLevBy/UrlShortener/internal/lib/logger/sl"
	"github.com/RomanLevBy/UrlShortener/internal/storage/sqlite"
	"log/slog"
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
	log.Info("Config init of", slog.String("env", conf.Env))
	log.Debug("Debug messages are enable")

	storage, err := sqlite.New(conf.StoragePath)
	if err != nil {
		log.Error("Failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	_ = storage

	//todo init router: chi, render

	//todo run server
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
