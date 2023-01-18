package setup

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"mentat-backend/internal/api"
	"mentat-backend/internal/config"
	"os"
)

func InitializeEcho() *echo.Echo {
	e := echo.New()
	api.SetupRoutes(e)

	env := os.Getenv("ENV")
	cfg, err := config.LoadConfig(config.DefaultConfigPath, env)
	if err != nil {
		log.Printf("could not load config from the beginning: %v", err)
		return nil
	}
	if cfg.Environment == "" {
		log.Printf("could not load config correctly: %v", err)
		return nil
	}

	ConfigHandler(e, cfg)

	return e
}
