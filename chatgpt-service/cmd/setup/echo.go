package setup

import (
	"chatgpt-service/internal/api"
	"chatgpt-service/internal/config"
	"chatgpt-service/internal/pkg/client"
	"github.com/labstack/echo/v4"
)

func InitializeEcho(cfg *config.GlobalConfig, oc client.OpenAIClient, fc client.FlipsideClient) (error, *echo.Echo) {
	e := echo.New()
	//ConfigHandler(e, *cfg, oc)
	err := api.SetupRoutes(e, *cfg, oc, fc)
	if err != nil {
		return err, nil
	}

	return nil, e
}
