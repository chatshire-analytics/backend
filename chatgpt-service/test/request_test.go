package test

import (
	"chatgpt-service/cmd/setup"
	"chatgpt-service/internal/api"
	"chatgpt-service/internal/config"
	cpkg "chatgpt-service/internal/pkg/client"
	"chatgpt-service/pkg/client"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTest(t *testing.T) (error, echo.Context) {
	cfg, err := config.LoadConfig(config.TestConfigPath, "dev")
	if err != nil {
		t.Errorf("could not load config: %v", err)
	}
	oc, err := setup.NewOpenAIClient(cfg)
	if err != nil {
		t.Errorf("could not create openai client: %v", err)
	}
	e := echo.New()
	e.GET(client.GetAllModels, api.ListModels)
	req := e.NewContext(httptest.NewRequest(echo.GET, client.GetAllModels, nil), httptest.NewRecorder())
	req.Set(cpkg.OpenAIClientKey, oc)
	return err, req
}

func TestListModels(t *testing.T) {
	// given
	err, req := setupTest(t)

	// when
	err = api.ListModels(req)

	// then
	if err != nil {
		t.Errorf("could not list models: %v", err)
	}
	res := req.Response()

	// then
	if res.Status != http.StatusOK {
		t.Errorf("expected status OK but got %v", res.Status)
	}
	body := res.Writer.(*httptest.ResponseRecorder).Body
	var listModelsResponse client.ListModelsResponse
	if err = json.Unmarshal(body.Bytes(), &listModelsResponse); err != nil {
		t.Errorf("could not unmarshal response: %v", err)
	}
	if len(listModelsResponse.Data) == 0 {
		t.Errorf("expected at least one model but got %v", len(listModelsResponse.Data))
	}
}
