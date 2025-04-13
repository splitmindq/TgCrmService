package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
)

type Bot struct {
	bot    *tgbotapi.BotAPI
	chatID int64
	logger *slog.Logger
}

func NewBot(token string, chatID int64, logger *slog.Logger) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		logger.Error("Failed to initialize Telegram bot", "error", err)
		return nil, err
	}
	logger.Info("Bot initialized")
	return &Bot{
		bot:    bot,
		chatID: chatID,
		logger: logger,
	}, nil
}

func (bot *Bot) Start() {
	bot.logger.Info("Starting Telegram bot")
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil || !update.Message.IsCommand() {
			continue
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		switch update.Message.Command() {
		case "start":
			msg.Text = "Привет! Я бот для лидов."
			bot.logger.Info("Received /start command", "chat_id", update.Message.Chat.ID)
		case "stats":
			msg.Text = "Статистика пока недоступна."
			bot.logger.Info("Received /stats command", "chat_id", update.Message.Chat.ID)
		default:
			msg.Text = "Неизвестная команда."
			bot.logger.Warn("Unknown command", "command", update.Message.Command())
		}

		if _, err := bot.bot.Send(msg); err != nil {
			bot.logger.Error("Failed to send response", "error", err)
		}
	}

}
