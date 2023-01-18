package api

import (
	"github.com/labstack/echo/v4"
	"mentat-backend/gateway-service/internal/pkg/telegram"
)

func SetupRoutes(e *echo.Echo) *echo.Echo {
	e.GET("/health", HealthCheck)
	e.POST("/telegram/webhook", telegram.WebhookHandler)
	return e
}
