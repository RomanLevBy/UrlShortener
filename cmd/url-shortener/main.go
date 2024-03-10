package main

import (
	//"fmt"
	"github.com/RomanLevBy/UrlShortener/internal/config"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	//export CONFIG_PATH=./config/local.yaml
	conf := config.MustLoad()

	log := setupLogger(conf.Env)
	log.Info("Config init of", slog.String("env", conf.Env))
	log.Debug("Debug messages are enable")

	//todo init storage

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
