package api

import (
	"chatgpt-service/internal/config"
	gclient "chatgpt-service/internal/pkg/client"
	"chatgpt-service/pkg/client"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"time"
)

func HttpRequestLogHandler(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := c.Request()
		res := c.Response()
		start := time.Now()
		if err := next(c); err != nil {
			c.Error(err)
		}
		stop := time.Now()
		logrus.WithFields(logrus.Fields{
			"component": "api",
			"timestamp": time.Now().Format(time.RFC3339),
			"method":    req.Method,
			"remote":    req.RemoteAddr,
			"path":      req.URL.Path,
			"proto":     req.Proto,
			"status":    res.Status,
			"duration":  stop.Sub(start),
		}).Log(logrus.InfoLevel, "request completed")
		return nil
	}
}

func SetupRoutes(e *echo.Echo, cfg config.GlobalConfig, oc gclient.OpenAIClient, fc gclient.FlipsideClient) error {
	hd, err := NewHandler(e.AcquireContext(), cfg, &oc, &fc)
	if err != nil {
		return err
	}
	e.Use(HttpRequestLogHandler)
	e.GET(HealthEndpoint, HealthCheck)
	e.GET(client.GetAllModels, hd.ListModels)
	e.GET(client.RetrieveModels, hd.RetrieveModel)
	e.POST(client.CreateCompletionEndpoint, hd.CreateCompletion)
	e.POST(gclient.CreateFlipsideQueryEndpoint, hd.CreateFlipsideQuery)

	return nil
}
