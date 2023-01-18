package api

import (
	"github.com/labstack/echo/v4"
	"mentat-backend/internal/api/telegram"
)

func SetupRoutes(e *echo.Echo) {
	e.GET("/health", HealthCheck)
	e.POST("/telegram/webhook", telegram.WebhookHandler)
}
