package setup

import (
	"chatgpt-service/internal/api"
	"chatgpt-service/internal/config"
	"chatgpt-service/internal/pkg/client"
	"github.com/labstack/echo/v4"
)

func InitializeEcho(cfg *config.GlobalConfig, oc *client.OpenAIClientInterface) (error, *echo.Echo) {
	e := echo.New()
	api.SetupRoutes(e)
	ConfigHandler(e, *cfg, oc)

	return nil, e
}
