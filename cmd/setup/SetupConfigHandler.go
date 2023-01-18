package setup

import (
	"github.com/labstack/echo/v4"
	"mentat-backend/internal/config"
)

func ConfigHandler(e *echo.Echo, cfg *config.GlobalConfig) {
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			ctx.Set(config.GlobalConfigKey, cfg)
			return next(ctx)
		}
	})
}
