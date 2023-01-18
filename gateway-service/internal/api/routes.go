package api

import "github.com/labstack/echo/v4"

func SetupRoutes(e *echo.Echo) *echo.Echo {
	e.GET("/health", HealthCheck)
	return e
}
