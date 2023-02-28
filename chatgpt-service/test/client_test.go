package test

import (
	"bufio"
	"bytes"
	"chatgpt-service/cmd/setup"
	"chatgpt-service/internal/api"
	"chatgpt-service/internal/config"
	cpkg "chatgpt-service/internal/pkg/client"
	"chatgpt-service/internal/pkg/constants"
	"chatgpt-service/pkg/client"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTest(t *testing.T, method string, endpoint string, bodyRaw *[]byte) (error, echo.Context, *api.Handler) {
	cfg, err := config.LoadConfig(config.TestConfigPath, "dev")
	if err != nil {
		t.Errorf("could not load config: %v", err)
	}
	oc, err := setup.NewOpenAIClient(cfg)
	if err != nil {
		t.Errorf("could not create openai client: %v", err)
	}
	fc, err := setup.NewFlipsideClient(cfg)
	if err != nil {
		t.Errorf("could not create flipside client: %v", err)
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
	ectx := e.NewContext(reqRaw, httptest.NewRecorder())
	ectx.Set(cpkg.OpenAIClientKey, oc)
	hd, err := api.NewHandler(ectx, *cfg, oc, fc)
	if err != nil {
		t.Errorf("could not create handler: %v", err)
	}
	return nil, ectx, hd
}

func setupTestSSE(t *testing.T, method string, endpoint string, bodyRaw *[]byte) (error, echo.Context, *api.Handler) {
	cfg, err := config.LoadConfig(config.TestConfigPath, "dev")
	if err != nil {
		t.Errorf("could not load config: %v", err)
	}
	oc, err := setup.NewOpenAIClient(cfg)
	if err != nil {
		t.Errorf("could not create openai client: %v", err)
	}
	fc, err := setup.NewFlipsideClient(cfg)
	if err != nil {
		t.Errorf("could not create flipside client: %v", err)
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
	ectx := e.NewContext(reqRaw, httptest.NewRecorder())
	ectx.Set(cpkg.OpenAIClientKey, oc)
	hd, err := api.NewHandler(ectx, *cfg, oc, fc)
	if err != nil {
		t.Errorf("could not create handler: %v", err)
	}
	return nil, ectx, hd
}

func TestListModels(t *testing.T) {
	// given
	err, ectx, hd := setupTest(t, http.MethodGet, client.GetAllModels, nil)
	if err != nil {
		t.Fatalf("could not create handler: %v", err)
	}

	// when
	err = hd.ListModels(ectx)

	// then
	if err != nil {
		t.Fatalf("could not list models: %v", err)
	}
	res := ectx.Response()
	if res.Status != http.StatusOK {
		t.Fatalf("expected status OK but got %v", res.Status)
	}
	body := res.Writer.(*httptest.ResponseRecorder).Body
	var listModelsResponse client.ListModelsResponse
	if err = json.Unmarshal(body.Bytes(), &listModelsResponse); err != nil {
		t.Fatalf("could not unmarshal response: %v", err)
	}
	if len(listModelsResponse.Data) == 0 {
		t.Fatalf("expected at least one model but got %v", len(listModelsResponse.Data))
	}
}

func TestRetrieveModel(t *testing.T) {
	// given
	err, ectx, hd := setupTest(t, http.MethodGet, client.RetrieveModels, nil)
	EXAMPLE_MODEL_ID := constants.TextDavinci003Engine
	ectx.SetParamNames(client.ModelIdParamKey)
	ectx.SetParamValues(EXAMPLE_MODEL_ID)

	// when
	err = hd.RetrieveModel(ectx)

	// then
	if err != nil {
		t.Errorf("could not retrieve model: %v", err)
	}
	res := ectx.Response()
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
	err, ectx, hd := setupTest(t, http.MethodPost, client.CreateCompletionEndpoint, &bodyRaw)
	if err != nil {
		t.Errorf("could not setup test: %v", err)
	}

	// when
	err = hd.CreateCompletion(ectx)

	// then
	if err != nil {
		t.Errorf("could not create completion: %v", err)
	}
	res := ectx.Response()
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
	err, ectx, hd := setupTestSSE(t, http.MethodPost, client.CreateCompletionEndpoint, &bodyRaw)
	if err != nil {
		t.Errorf("could not setup test: %v", err)
	}

	// when
	err = hd.CreateCompletionStream(ectx)

	// then
	if err != nil {
		t.Errorf("could not create completion: %v", err)
	}
	res := ectx.Response()
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
			if err == io.EOF {
				break
			}
		}
		rawString += string(line)
	}
	fmt.Println("rawString", rawString)
	assert.NotEmpty(t, rawString)
}

func TestFlipsideCryptoCreateAQuery(t *testing.T) {
	// given
	bodyTest := &client.CreateFlipsideQueryRequest{
		Sql:        "select TX_ID from ethereum.TRANSACTIONS limit 10",
		TtlMinutes: 15,
		Cache:      true,
		Params: struct {
			AdditionalProp1 string `json:"additionalProp1"`
			AdditionalProp2 string `json:"additionalProp2"`
			AdditionalProp3 string `json:"additionalProp3"`
		}{
			AdditionalProp1: "string",
			AdditionalProp2: "string",
			AdditionalProp3: "string",
		},
	}

	bodyRaw, err := json.Marshal(bodyTest)
	if err != nil {
		t.Errorf("could not marshal request body: %v", err)
	}

	err, ectx, hd := setupTest(t, http.MethodGet, client.GetAllModels, &bodyRaw)
	if err != nil {
		t.Fatalf("could not create handler: %v", err)
	}

	// when
	err = hd.CreateFlipsideQuery(ectx)

	// then
	if err != nil {
		t.Fatalf("could not create query: %v", err)
	}
	res := ectx.Response()
	if res.Status != http.StatusOK {
		t.Fatalf("expected status OK but got %v", res.Status)
	}
	bodyVerify := res.Writer.(*httptest.ResponseRecorder).Body
	var queryResponse client.GetFlipsideQueryResultSuccessResponse
	if err = json.Unmarshal(bodyVerify.Bytes(), &queryResponse); err != nil {
		t.Fatalf("could not unmarshal response: %v", err)
	}

	fmt.Println("queryResponse", queryResponse)
}

func TestE2E_1(t *testing.T) {
	// given
	bodyTest := &client.CreateFlipsideQueryRequest{
		Sql:        "select TX_ID from ethereum.TRANSACTIONS limit 10",
		TtlMinutes: 15,
		Cache:      true,
		Params: struct {
			AdditionalProp1 string `json:"additionalProp1"`
			AdditionalProp2 string `json:"additionalProp2"`
			AdditionalProp3 string `json:"additionalProp3"`
		}{
			AdditionalProp1: "string",
			AdditionalProp2: "string",
			AdditionalProp3: "string",
		},
	}

	bodyRaw, err := json.Marshal(bodyTest)
	if err != nil {
		t.Errorf("could not marshal request body: %v", err)
	}

	err, ectx, hd := setupTest(t, http.MethodGet, client.GetAllModels, &bodyRaw)
	if err != nil {
		t.Fatalf("could not create handler: %v", err)
	}

	// when
	err = hd.CreateFlipsideQuery(ectx)

	// then
	if err != nil {
		t.Fatalf("could not create query: %v", err)
	}
	res := ectx.Response()
	if res.Status != http.StatusOK {
		t.Fatalf("expected status OK but got %v", res.Status)
	}
	bodyVerify := res.Writer.(*httptest.ResponseRecorder).Body
	var queryResponse client.GetFlipsideQueryResultSuccessResponse
	if err = json.Unmarshal(bodyVerify.Bytes(), &queryResponse); err != nil {
		t.Fatalf("could not unmarshal response: %v", err)
	}

	fmt.Println("queryResponse", queryResponse)
}

func TestFlipsideCryptoGetQueryResult(t *testing.T) {
	//
}

func TestChatGptToFlipsideCryptoExample(t *testing.T) {
	//
}
