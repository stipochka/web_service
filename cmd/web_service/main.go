package main

import (
	"log/slog"
	"os"

	"github.com/stipochka/web_service/internal/config"
	handler "github.com/stipochka/web_service/internal/handlers"
	"github.com/stipochka/web_service/internal/repository"
	"github.com/stipochka/web_service/internal/service"
)

func main() {
	cfg := config.MustLoad()

	logger := setupLogger(cfg.Env)

	clientConn, err := repository.NewGRPCConn(cfg.GRPCServer.Address)

	if err != nil {
		logger.Error("failed to create conn", slog.String("error", err.Error()))
		return
	}

	repo := repository.NewGRPCRepository(clientConn)

	serviceLayer := service.NewService(repo)

	h := handler.NewHandler(logger, serviceLayer)

	router := h.InitRoutes()

	logger.Info("Server is running on address", slog.String("address", cfg.HTTPServer.Address))
	if err := router.Run(cfg.HTTPServer.Address); err != nil {
		logger.Error("Error while running server", slog.Any("error", err))
		panic(err)
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "local":
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case "dev":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	case "prod":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
