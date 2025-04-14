package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
)

type Bot struct {
	Bot    *tgbotapi.BotAPI
	ChatID int64
	Logger *slog.Logger
}

func NewBot(token string, chatID int64, logger *slog.Logger) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		logger.Error("Failed to initialize Telegram bot", "error", err)
		return nil, err
	}
	logger.Info("Bot initialized")
	return &Bot{
		Bot:    bot,
		ChatID: chatID,
		Logger: logger,
	}, nil
}

func (bot *Bot) Start() {
	bot.Logger.Info("Starting Telegram bot")
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.Bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil || !update.Message.IsCommand() {
			continue
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		switch update.Message.Command() {
		case "start":
			msg.Text = "Привет! Я бот для лидов."
			bot.Logger.Info("Received /start command", "chat_id", update.Message.Chat.ID)
		case "stats":
			msg.Text = "Статистика пока недоступна."
			bot.Logger.Info("Received /stats command", "chat_id", update.Message.Chat.ID)
		default:
			msg.Text = "Неизвестная команда."
			bot.Logger.Warn("Unknown command", "command", update.Message.Command())
		}

		if _, err := bot.Bot.Send(msg); err != nil {
			bot.Logger.Error("Failed to send response", "error", err)
		}
	}

}
func (b *Bot) SendNotification(message string) error {
	msg := tgbotapi.NewMessage(b.ChatID, message)
	_, err := b.Bot.Send(msg)
	if err != nil {
		b.Logger.Error("failed to send notification", "error", err)
		return fmt.Errorf("send notification: %w", err)
	}
	return nil
}
