package telegram

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	telegram2 "mentat-backend/pkg/api/telegram"
	"net/http"
	"os"
)

func WebhookHandler(c echo.Context) error {
	body := &telegram2.WebhookReqBody{}
	if err := c.Bind(body); err != nil {
		return c.String(http.StatusBadRequest, "could not decode request body: "+err.Error())
	}

	if err := sayPolo(body.Message.Chat.ID); err != nil {
		log.Printf("could not send message: %v", err)
		return c.String(http.StatusInternalServerError, "could not send message")
	}

	log.Printf("received message: %v", body.Message)
	return c.String(http.StatusOK, "ok")
}

func sayPolo(chatID int64) error {
	reqBody := &telegram2.WebhookResBody{
		ChatID: chatID,
		Text:   "Bitcoin to the moon!!",
	}
	// Create the JSON body from the struct
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	apiId := os.Getenv("DEV_TELEGRAM_TOKEN_API_ID")
	apiToken := os.Getenv("DEV_TELEGRAM_TOKEN_API_KEY")

	url := fmt.Sprintf("https://api.telegram.org/bot%s:%s/sendMessage", apiId, apiToken)
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
