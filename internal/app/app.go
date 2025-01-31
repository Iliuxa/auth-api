package app

import (
	grpcapp "auth-api/internal/app/grpc"
	"auth-api/internal/repository"
	"auth-api/internal/usecase"
	"database/sql"
	_ "github.com/lib/pq"
	"log/slog"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
	DB         *sql.DB
}

func New(log *slog.Logger, grpcPort int, dataSourceName string, tokenTTL time.Duration) *App {

	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		panic(err)
	}

	userRepository := repository.NewUserRepository(db)
	authUsecase := usecase.NewAuthUsecase(userRepository, log, tokenTTL)
	grpcApp := grpcapp.New(log, authUsecase, grpcPort)

	return &App{
		GRPCServer: grpcApp,
		DB:         db,
	}
}
