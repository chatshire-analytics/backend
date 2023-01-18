package setup

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
)

func SetupTelegramBot() error {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		return err
	}

	bot.Debug = true

	// The offset is set to 0 which means that the bot will start receiving updates from the most recent update.
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		if _, err := bot.Send(msg); err != nil {
			return err
		}
	}

	return nil
}
