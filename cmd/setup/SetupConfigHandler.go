package setup

import (
	"github.com/labstack/echo/v4"
	"mentat-backend/internal/config"
)

func SetConfigHandler(e *echo.Echo, cfg *config.GlobalConfig) *echo.Echo {
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			ctx.Set("config", cfg)
			return next(ctx)
		}
	})
	return e
}
