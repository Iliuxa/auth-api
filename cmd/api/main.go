package main

import (
	"auth-api/internal/config"
	"log/slog"
	"os"
)

func main() {
	cfg := config.MustLoad()

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	_ = cfg
	_ = log

	// TODO: инициализировать приложение (app)

	// TODO: запустить gRPC-сервер приложения

}
