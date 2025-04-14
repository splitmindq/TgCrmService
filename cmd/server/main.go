package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	//"github.com/go-chi/chi/v5/middleware"
	//"github.com/go-chi/cors"
	"lead-bitrix/internal/config"
	"lead-bitrix/internal/http-server/handlers/create"
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

	//todo init configs (database)

	cfg := config.MustLoadConfig()
	fmt.Println(cfg)

	logger := setupLogger(cfg.Env)

	bot, _ := telegram.NewBot(cfg.TelegramConfig.Token, cfg.TelegramConfig.ChatID, logger)
	go bot.Start()

	router := chi.NewRouter()
	router.Post("/save", create.NewLead(logger, bot))

	//todo init storage
	//todo init router

	logger.Info("listening server")
	err := cfg.HTTPListen(router)
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
