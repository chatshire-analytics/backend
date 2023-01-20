package test

import (
	"chatgpt-service/internal/api"
	"chatgpt-service/pkg/client"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"testing"
)

func TestListModels(t *testing.T) {
	// given
	e := echo.New()
	e.GET(client.GetAllModels, api.ListModels)
	req := httptest.NewRequest(echo.GET, client.GetAllModels, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// when
	if assert.NoError(t, api.ListModels(c)) {
		// then
		assert.Equal(t, 200, rec.Code)
	}
}
