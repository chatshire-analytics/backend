package setup

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
	"os"
)

func SetupTelegramBot(token string, webhookURL string) error {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return err
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	certFile := os.Getenv("CERT_PATH")
	file, err := os.ReadFile(certFile)
	if err != nil {
		return err
	}

	wh, err := tgbotapi.NewWebhookWithCert(webhookURL, tgbotapi.FileBytes{Name: "telegram-cert.pem", Bytes: file})
	if err != nil {
		log.Printf("Error creating webhook: %v", err)
		return err
	}

	_, err = bot.Request(wh)
	if err != nil {
		log.Printf("Error setting webhook: %v", err)
		return err
	}

	info, err := bot.GetWebhookInfo()
	if err != nil {
		log.Printf("Error getting webhook info: %v", err)
		return err
	}

	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
		return errors.New("telegram callback failed")
	}

	updates := bot.ListenForWebhook("/" + bot.Token)
	// set channel for error
	errChan := make(chan error, 1)
	go func() {
		err := http.ListenAndServeTLS("0.0.0.0:8443", "cert.pem", "key.pem", nil)
		if err != nil {
			log.Printf("Error starting webhook: %v", err)
			errChan <- err
		}
		close(errChan)
	}()

	for update := range updates {
		log.Printf("%+v\n", update)
	}

	select {
	case err := <-errChan:
		return err
	default:
		return nil
	}
}
