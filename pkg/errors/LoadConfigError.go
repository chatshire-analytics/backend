package setup

import (
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"mentat-backend/cmd/setup"
	"mentat-backend/internal/api"
	"mentat-backend/internal/config"
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
		return errors.New("failed to load config: environment string is empty"), nil
	}

	setup.ConfigHandler(e, cfg)

	return nil, e
}
