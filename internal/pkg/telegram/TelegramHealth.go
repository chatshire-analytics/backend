package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"mentat-backend/internal/config"
	tgpkg "mentat-backend/pkg/api/telegram"
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
		return err
	}

	cfg := c.Get(config.GlobalConfigKey).(*config.GlobalConfig)
	apiId := cfg.TelegramEnv.API_ID
	apiKey := cfg.TelegramEnv.API_KEY

	url := fmt.Sprintf("https://api.telegram.org/bot%s:%s/sendMessage", apiId, apiKey)
	fmt.Println("url: ", url)
	res, err := http.Post(url, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("unexpected status" + res.Status)
	}

	return nil
}
