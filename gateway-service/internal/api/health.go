package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type Health struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, Health{
		Name:        "Mentat Backend",
		Description: "Mentat Backend API",
	})
}
