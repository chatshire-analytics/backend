package setup

import (
	"github.com/labstack/echo/v4"
	"mentat-backend/gateway-service/internal/api"
)

func InitializeEcho() *echo.Echo {
	e := echo.New()
	e = api.SetupRoutes(e)
	return e
}
