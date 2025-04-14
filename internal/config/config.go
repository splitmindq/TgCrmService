package config

import (
	"github.com/go-chi/chi/v5"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"
)

type Config struct {
	Env            string         `yaml:"env" env-default:"local"`
	HTTPServer     HTTPServer     `yaml:"http_server" env-required:"true"`
	TelegramConfig TelegramConfig `yaml:"telegram_config" env-required:"true"`
}

type HTTPServer struct {
	Address     string        `yaml:"addr" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type TelegramConfig struct {
	Token  string `yaml:"token" env:"BOT_TOKEN" env-required:"true"`
	ChatID int64  `yaml:"chat_id" env:"CHAT_ID" env-required:"true"`
}

func MustLoadConfig() *Config {

	err := godotenv.Load()

	configPath := os.Getenv("CONFIG_PATH")

	if configPath == "" {
		log.Fatal("CONFIG_PATH environment variable not set")
	}

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatal("Config file does not exist")
	}

	var config Config

	err = cleanenv.ReadConfig(configPath, &config)

	if err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}

	return &config
}

func (config *Config) HTTPListen(router chi.Router) error {

	server := http.Server{
		Addr:         config.HTTPServer.Address,
		ReadTimeout:  config.HTTPServer.Timeout,
		WriteTimeout: config.HTTPServer.Timeout,
		IdleTimeout:  config.HTTPServer.IdleTimeout,
		Handler:      router,
	}

	return server.ListenAndServe()

}
