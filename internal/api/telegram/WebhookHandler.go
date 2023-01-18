package telegram

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"mentat-backend/internal/pkg/telegram"
	tgpkg "mentat-backend/pkg/api/telegram"
	"net/http"
)

func WebhookHandler(c echo.Context) error {
	body := &tgpkg.WebhookReqBody{}
	if err := c.Bind(body); err != nil {
		return c.String(http.StatusBadRequest, "could not decode request body: "+err.Error())
	}

	if err := telegram.WebhookHealth(c, body); err != nil {
		log.Printf("could not send message: %v", err)
		return c.String(http.StatusInternalServerError, "could not send message")
	}

	log.Printf("received message: %v", body.Message)
	return c.String(http.StatusOK, "ok")
}
