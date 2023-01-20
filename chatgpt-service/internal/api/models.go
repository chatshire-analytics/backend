package api

import (
	"chatgpt-service/internal/pkg/client"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

func ListModels(c echo.Context) error {
	ocInterface := c.Get(client.OpenAIClientKey)
	oc, ok := ocInterface.(*client.OpenAIClient)
	if !ok {
		return errors.New("could not convert to OpenAI client")
	}
	res, err := oc.ListModels(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(200, res)
}
