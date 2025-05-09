package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"lead-bitrix/internal/config"
	"lead-bitrix/internal/http-server/handlers/lead"
	"lead-bitrix/internal/storage/pgx"
	"lead-bitrix/internal/telegram"
	"log"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envProd  = "prod"
	envDev   = "dev"
)

func main() {

	cfg := config.MustLoadConfig()
	fmt.Println(cfg)

	logger := setupLogger(cfg.Env)

	storage, err := pgx.NewStorage(cfg)
	if err != nil {
		logger.Error(err.Error())
	}
	defer storage.Close()

	bot, _ := telegram.NewBot(cfg.TelegramConfig.Token, cfg.TelegramConfig.ChatID, logger)

	go bot.Start()

	router := chi.NewRouter()

	router.Use(middleware.RequestID)

	router.Get("/api/leads", lead.GetLeads(logger, storage))
	router.Get("/api/leads/{email}", lead.LeadGetByEmail(logger, storage))
	router.Post("/api/leads", lead.NewLead(logger, bot, storage))
	router.Delete("/api/leads/{email}", lead.DelLead(logger, bot, storage))
	router.Patch("/api/leads/{phone}", lead.UpdateLead(logger, storage))

	logger.Info("listening server")
	err = cfg.HTTPListen(router)
	if err != nil {
		log.Fatalf("Could not start server: %v", err)
	}

	return

}

func setupLogger(env string) *slog.Logger {

	var logger *slog.Logger

	switch env {
	case envLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return logger
}
