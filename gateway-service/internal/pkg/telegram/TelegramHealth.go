package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"mentat-backend/internal/config"
	tgpkg "mentat-backend/pkg/api/telegram"
	setup "mentat-backend/pkg/errors"
	"net/http"
)

func WebhookHealth(c echo.Context, body *tgpkg.WebhookReqBody) error {
	chatId := body.Message.Chat.ID
	reqBody := &tgpkg.WebhookResBody{
		ChatId: chatId,
		Text:   "Bitcoin to the moon!!",
	}
	// Create the JSON body from the struct
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"component": "telegram",
			"ChatId":    chatId,
			"Text":      reqBody.Text,
		}).WithError(err).Error(logrus.ErrorLevel, "could not marshal request body")
		return err
	}

	cfg := c.Get(config.GlobalConfigKey).(*config.GlobalConfig)
	apiId := cfg.TelegramEnv.API_ID
	apiKey := cfg.TelegramEnv.API_KEY

	url := fmt.Sprintf("https://api.telegram.org/bot%s:%s/sendMessage", apiId, apiKey)
	res, err := http.Post(url, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"component": "telegram",
			"ChatId":    chatId,
			"Text":      reqBody.Text,
			"URL":       url,
			"API_ID":    apiId,
			"API_KEY":   apiKey,
		}).WithError(err).Error(logrus.ErrorLevel, "could not send message")
		return err
	}
	if res.StatusCode != http.StatusOK {
		logrus.WithFields(logrus.Fields{
			"component": "telegram",
			"ChatId":    chatId,
			"Text":      reqBody.Text,
			"URL":       url,
			"API_ID":    apiId,
			"API_KEY":   apiKey,
			"Status":    res.Status,
		}).WithError(err).Error(logrus.ErrorLevel, "unexpected status"+res.Status)
		return setup.TelegramError(res.Status)
	}

	return nil
}
