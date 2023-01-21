package test

import (
	"bufio"
	"bytes"
	"chatgpt-service/cmd/setup"
	"chatgpt-service/internal/api"
	"chatgpt-service/internal/config"
	cpkg "chatgpt-service/internal/pkg/client"
	"chatgpt-service/internal/pkg/engine"
	"chatgpt-service/pkg/client"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupTest(t *testing.T, method string, endpoint string, bodyRaw *[]byte) (error, echo.Context) {
	cfg, err := config.LoadConfig(config.TestConfigPath, "dev")
	if err != nil {
		t.Errorf("could not load config: %v", err)
	}
	oc, err := setup.NewOpenAIClient(cfg)
	if err != nil {
		t.Errorf("could not create openai client: %v", err)
	}
	e := echo.New()
	var body io.Reader
	if bodyRaw == nil {
		body = nil
	} else {
		body = bytes.NewBuffer(*bodyRaw)
	}
	reqRaw := httptest.NewRequest(method, endpoint, body)
	reqRaw.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req := e.NewContext(reqRaw, httptest.NewRecorder())
	req.Set(cpkg.OpenAIClientKey, oc)
	return nil, req
}

func setupTestSSE(t *testing.T, method string, endpoint string, bodyRaw *[]byte) (error, echo.Context) {
	cfg, err := config.LoadConfig(config.TestConfigPath, "dev")
	if err != nil {
		t.Errorf("could not load config: %v", err)
	}
	oc, err := setup.NewOpenAIClient(cfg)
	if err != nil {
		t.Errorf("could not create openai client: %v", err)
	}
	e := echo.New()
	var body io.Reader
	if bodyRaw == nil {
		body = nil
	} else {
		body = bytes.NewBuffer(*bodyRaw)
	}
	reqRaw := httptest.NewRequest(method, endpoint, body)
	reqRaw.Header.Set(echo.HeaderCacheControl, "no-cache")
	reqRaw.Header.Set(echo.HeaderAccept, "text/event-stream")
	reqRaw.Header.Set(echo.HeaderConnection, "keep-alive")
	reqRaw.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req := e.NewContext(reqRaw, httptest.NewRecorder())
	req.Set(cpkg.OpenAIClientKey, oc)
	return nil, req
}

func TestListModels(t *testing.T) {
	// given
	err, req := setupTest(t, http.MethodGet, client.GetAllModels, nil)

	// when
	err = api.ListModels(req)

	// then
	if err != nil {
		t.Errorf("could not list models: %v", err)
	}
	res := req.Response()
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

func TestRetrieveModel(t *testing.T) {
	// given
	err, req := setupTest(t, http.MethodGet, client.RetrieveModels, nil)
	MODEL_ID_KEY := "model_id"
	EXAMPLE_MODEL_ID := engine.TextDavinci003Engine
	req.SetParamNames(MODEL_ID_KEY)
	req.SetParamValues(EXAMPLE_MODEL_ID)

	// when
	err = api.RetrieveModel(req)

	// then
	if err != nil {
		t.Errorf("could not retrieve model: %v", err)
	}
	res := req.Response()
	if res.Status != http.StatusOK {
		t.Errorf("expected status OK but got %v", res.Status)
	}
	body := res.Writer.(*httptest.ResponseRecorder).Body
	var retrievedModelObject client.ModelObject
	if err = json.Unmarshal(body.Bytes(), &retrievedModelObject); err != nil {
		t.Errorf("could not unmarshal response: %v", err)
	}
	if retrievedModelObject.ID != "text-davinci-003" {
		t.Errorf("expected model with id text-davinci-003 but got %v", retrievedModelObject.ID)
	}
}

func TestCreateCompletion(t *testing.T) {
	// given
	bodyTest := client.NewCompletionRequest("this is a test", 3, nil, nil)
	bodyRaw, err := json.Marshal(bodyTest)
	if err != nil {
		t.Errorf("could not marshal request body: %v", err)
	}
	err, req := setupTest(t, http.MethodPost, client.CreateCompletionEndpoint, &bodyRaw)
	if err != nil {
		t.Errorf("could not setup test: %v", err)
	}

	// when
	err = api.CreateCompletion(req)

	// then
	if err != nil {
		t.Errorf("could not create completion: %v", err)
	}
	res := req.Response()
	if res.Status != http.StatusOK {
		t.Errorf("expected status OK but got %v", res.Status)
	}
	bodyVerify := res.Writer.(*httptest.ResponseRecorder).Body
	var completionResponse client.CompletionResponse
	if err = json.Unmarshal(bodyVerify.Bytes(), &completionResponse); err != nil {
		t.Errorf("could not unmarshal response: %v", err)
	}
	if len(completionResponse.Choices) == 0 {
		t.Errorf("expected at least one completion but got %v", len(completionResponse.Choices))
	}
}

func TestCreateCompletionStreamTrue(t *testing.T) {
	// given
	stream := new(bool)
	*stream = true
	bodyTest := client.NewCompletionRequest("one thing that you should know about golang", 20, nil, stream)
	bodyRaw, err := json.Marshal(bodyTest)
	if err != nil {
		t.Errorf("could not marshal request body: %v", err)
	}
	err, req := setupTestSSE(t, http.MethodPost, client.CreateCompletionEndpoint, &bodyRaw)
	if err != nil {
		t.Errorf("could not setup test: %v", err)
	}

	// when
	err = api.CreateCompletionStream(req)

	// then
	if err != nil {
		t.Errorf("could not create completion: %v", err)
	}
	res := req.Response()
	if res.Status != http.StatusOK {
		t.Errorf("expected status OK but got %v", res.Status)
	}
	if res.Header().Get(echo.HeaderContentType) != "text/event-stream" {
		t.Errorf("expected content type text/event-stream but got %v", res.Header().Get(echo.HeaderContentType))
	}
	var rawString string
	reader := bufio.NewReader(res.Writer.(*httptest.ResponseRecorder).Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if strings.HasSuffix(string(line), "\\end ") {
				break
			}
		}
		rawString += string(line)
	}
	fmt.Println("rawString", rawString)
	assert.NotEmpty(t, rawString)
}
