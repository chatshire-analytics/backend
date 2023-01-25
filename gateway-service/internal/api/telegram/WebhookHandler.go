package telegram

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
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
		logrus.WithFields(logrus.Fields{
			"component": "telegram",
		}).WithError(err).Error(logrus.ErrorLevel, "could not send message")
		return c.String(http.StatusInternalServerError, "could not send message to the client: "+err.Error())
	}

	logrus.WithFields(logrus.Fields{
		"component": "telegram",
	}).Log(logrus.InfoLevel, "received message: %v", body.Message)

	return c.String(http.StatusOK, "ok")
}
