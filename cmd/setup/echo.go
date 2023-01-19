package setup

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"mentat-backend/internal/api"
	"mentat-backend/internal/config"
	setup "mentat-backend/pkg/errors"
	"os"
)

func InitializeEcho() (error, *echo.Echo) {
	e := echo.New()
	api.SetupRoutes(e)

	env := os.Getenv("ENV")
	logrus.WithFields(logrus.Fields{
		"component": "setup",
		"env":       env,
	}).Log(logrus.InfoLevel, "ENV is set to "+env)

	cfg, err := config.LoadConfig(config.DefaultConfigPath, env)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"component": "setup",
			"env":       env,
		}).WithError(err).Error(logrus.ErrorLevel, "Failed to load config")
		return err, nil
	}
	if cfg.Environment == "" {
		logrus.WithFields(logrus.Fields{
			"component": "setup",
			"env":       env,
		}).WithError(err).Error(logrus.ErrorLevel, "Failed to load config")
		return setup.LoadConfigError(), nil
	}

	ConfigHandler(e, cfg)

	return nil, e
}
