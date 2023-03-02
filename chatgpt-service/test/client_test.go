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
	"time"
)

func setupTest(t *testing.T, method string, endpoint string, bodyRaw *[]byte, paramStr *string) (error, echo.Context, *api.Handler) {
	cfg, err := config.LoadConfig(config.TestConfigPath, "dev")
	if err != nil {
		t.Errorf("could not load config: %v", err)
	}
	db := setup.InitializeDatabase(cfg)
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
	if bodyRaw != nil {
		reqRaw := httptest.NewRequest(method, endpoint, body)
		reqRaw.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		ectx := e.NewContext(reqRaw, httptest.NewRecorder())
		ectx.Set(cpkg.OpenAIClientKey, oc)
		ectx.Set(cpkg.FlipsideClientKey, fc)
		hd, err := api.NewHandler(ectx, *cfg, oc, fc, db)
		if err != nil {
			t.Errorf("could not create handler: %v", err)
		}
		return nil, ectx, hd
	}
	// path parameter BindPathParams
	if paramStr != nil {
		reqRaw := httptest.NewRequest(method, endpoint, nil)
		reqRaw.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		ectx := e.NewContext(reqRaw, httptest.NewRecorder())
		ectx.SetParamNames("token")
		ectx.SetParamValues(*paramStr)
		ectx.Set(cpkg.OpenAIClientKey, oc)
		ectx.Set(cpkg.FlipsideClientKey, fc)
		hd, err := api.NewHandler(ectx, *cfg, oc, fc, db)
		if err != nil {
			t.Errorf("could not create handler: %v", err)
		}
		return nil, ectx, hd
	}
	return fmt.Errorf("could not create new http test handler"), nil, nil
}

func setupTestSSE(t *testing.T, method string, endpoint string, bodyRaw *[]byte) (error, echo.Context, *api.Handler) {
	cfg, err := config.LoadConfig(config.TestConfigPath, "dev")
	if err != nil {
		t.Errorf("could not load config: %v", err)
	}
	db := setup.InitializeDatabase(cfg)
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
	hd, err := api.NewHandler(ectx, *cfg, oc, fc, db)
	if err != nil {
		t.Errorf("could not create handler: %v", err)
	}
	return nil, ectx, hd
}

func TestListModels(t *testing.T) {
	// given
	err, ectx, hd := setupTest(t, http.MethodGet, client.GetAllModels, nil, nil)
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
	err, ectx, hd := setupTest(t, http.MethodGet, client.RetrieveModels, nil, nil)
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
	err, ectx, hd := setupTest(t, http.MethodPost, client.CreateCompletionEndpoint, &bodyRaw, nil)
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

	err, ectx, hd := setupTest(t, http.MethodGet, client.GetAllModels, &bodyRaw, nil)
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
	var queryResponse client.CreateFlipsideQuerySuccessResponse
	if err = json.Unmarshal(bodyVerify.Bytes(), &queryResponse); err != nil {
		t.Fatalf("could not unmarshal response: %v", err)
	}

	fmt.Println("queryResponse", queryResponse)
}

func TestFlipsideCryptoGetQueryResult(t *testing.T) {
	// given
	paramTest := &client.GetFlipsideQueryResultRequest{
		Token: "queryRun-eccac5ecb5ca3814d85de2557318359e",
		//Token: "queryRun-b4f7626374751285a9eb32bf477ec2ee",
		//Token: "queryRun-8d8c035e974d12bb180f8f8dd898b1ba",
	}

	err, ectx, hd := setupTest(t, http.MethodPost, cpkg.GetFlipsideQueryResultEndpoint, nil, &paramTest.Token)
	if err != nil {
		t.Fatalf("could not create handler: %v", err)
	}

	// when
	err = hd.GetFlipsideQueryResult(ectx)
	if err != nil {
		t.Fatalf("could not get query result: %v", err)
	}

	// then
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

func TestFlipsideCryptoGetQueryRunningResult(t *testing.T) {
	// given
	paramTest := &client.GetFlipsideQueryResultRequest{
		//Token: "queryRun-b4f7626374751285a9eb32bf477ec2ee",
		Token: "queryRun-8d8c035e974d12bb180f8f8dd898b1ba",
	}

	err, ectx, hd := setupTest(t, http.MethodPost, cpkg.GetFlipsideQueryResultEndpoint, nil, &paramTest.Token)
	if err != nil {
		t.Fatalf("could not create handler: %v", err)
	}

	// when
	err = hd.GetFlipsideQueryResult(ectx)
	if err != nil {
		t.Fatalf("could not get query result: %v", err)
	}

	// then
	res := ectx.Response()
	if res.Status != http.StatusAccepted {
		t.Fatalf("expected status OK but got %v", res.Status)
	}
}

// INTEGRATION TEST !!!

// // NOTE: PLEASE CONFIGURE PYTHON PATH BEFORE RUNNING THE TEST!
// THIS SHOULD BE HANDLED IN FURTHER VERSIONS
//	// python -c "import sys; print(sys.prefix)"
//	// PATH=/Users/sigridjin.eth/Documents/github/backend/venv/bin:$PATH
//
//	//cmd2 := exec.Command("python", "-c", "import sys; print(sys.prefix)")
//	//byte, err := cmd2.Output()
//	//byte = bytes.TrimSpace(byte)
//	//fmt.Println(string(byte))
//	//if err != nil {
//	//	t.Fatal(err)
//	//}
//
//	//installCmd := exec.Command("pip", "install", "-r", "requirements.txt")
//	//err = installCmd.Run()
//	//if err != nil {
//	//	t.Fatal(err)
//	//}

func TestIntegration_1_ReversedPython(t *testing.T) {
	// STAGE 1. GPT
	//schemaRaw, err := engine.CreatePrompt("Find the most expensive gas fees-cost transaction in last 7 days.")
	bodyInPrompt := client.GPTPromptRequest{
		Prompt: "Find the most expensive gas fees-cost transaction in last 7 days.",
	}
	bodyInPromptRaw, err := json.Marshal(bodyInPrompt)
	if err != nil {
		t.Errorf("could not marshal request body: %v", err)
	}
	err_gpt, ectx_gpt, hd_gpt := setupTest(t, http.MethodGet, cpkg.GPTGenerateQueryEndpoint, &bodyInPromptRaw, nil)
	if err_gpt != nil {
		t.Fatalf("could not create handler: %v", err_gpt)
	}
	errGpt2 := hd_gpt.RunGptPythonClient(ectx_gpt)
	if errGpt2 != nil {
		t.Fatalf("could not create query: %v", errGpt2)
	}
	res_gpt := ectx_gpt.Response()
	if res_gpt.Status != http.StatusOK {
		t.Fatalf("expected status OK but got %v", res_gpt.Status)
	}
	bodyverifyGpt := res_gpt.Writer.(*httptest.ResponseRecorder).Body
	//var gptResponse client.GPTPromptSuccessfulResponse
	//if err_gpt = json.Unmarshal(bodyverifyGpt.Bytes(), &gptResponse); err_gpt != nil {
	//	t.Fatalf("could not unmarshal response: %v", err_gpt)
	//}
	gptResponse := string(bodyverifyGpt.Bytes())
	fmt.Println("gptResponse", gptResponse)

	// 2. Create a query
	bodyTest := &client.CreateFlipsideQueryRequest{
		Sql:        gptResponse,
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

	err, ectx, hd := setupTest(t, http.MethodGet, client.GetAllModels, &bodyRaw, nil)
	if err != nil {
		t.Fatalf("could not create handler: %v", err)
	}
	err = hd.CreateFlipsideQuery(ectx)
	if err != nil {
		t.Fatalf("could not create query: %v", err)
	}
	res := ectx.Response()
	if res.Status != http.StatusOK {
		t.Fatalf("expected status OK but got %v", res.Status)
	}
	bodyVerify := res.Writer.(*httptest.ResponseRecorder).Body
	var queryResponse client.CreateFlipsideQuerySuccessResponse
	if err = json.Unmarshal(bodyVerify.Bytes(), &queryResponse); err != nil {
		t.Fatalf("could not unmarshal response: %v", err)
	}

	fmt.Println("queryResponse", queryResponse)
}

func TestEndToEnd_1_ReversedPython(t *testing.T) {
	bodyInPrompt := client.GPTPromptRequest{
		Prompt: "Find the most expensive gas fees-cost transaction in last 7 days.",
	}
	bodyInPromptRaw, err := json.Marshal(bodyInPrompt)
	if err != nil {
		t.Errorf("could not marshal request body: %v", err)
	}
	err_gpt, ectx_gpt, hd_gpt := setupTest(t, http.MethodGet, cpkg.GPTGenerateQueryEndpoint, &bodyInPromptRaw, nil)
	if err_gpt != nil {
		t.Fatalf("could not create handler: %v", err_gpt)
	}
	errGpt2 := hd_gpt.RunGptPythonClient(ectx_gpt)
	if errGpt2 != nil {
		t.Fatalf("could not create query: %v", errGpt2)
	}
	res_gpt := ectx_gpt.Response()
	if res_gpt.Status != http.StatusOK {
		t.Fatalf("expected status OK but got %v", res_gpt.Status)
	}
	fmt.Println("res_gpt", res_gpt)

	bodyverifyGpt := res_gpt.Writer.(*httptest.ResponseRecorder).Body
	var gptQueryResponse client.GPTPromptSuccessfulResponse
	if err_gpt = json.Unmarshal(bodyverifyGpt.Bytes(), &gptQueryResponse); err_gpt != nil {
		t.Fatalf("could not unmarshal response: %v", err_gpt)
	}

	bodyTest := &client.CreateFlipsideQueryRequest{
		Id:         string(rune(gptQueryResponse.Id)),
		Sql:        gptQueryResponse.Result,
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

	err, ectx, hd := setupTest(t, http.MethodGet, client.GetAllModels, &bodyRaw, nil)
	if err != nil {
		t.Fatalf("could not create handler: %v", err)
	}

	err = hd.CreateFlipsideQuery(ectx)
	if err != nil {
		t.Fatalf("could not create query: %v", err)
	}

	res := ectx.Response()
	if res.Status != http.StatusOK {
		t.Fatalf("expected status OK but got %v", res.Status)
	}

	bodyVerify := res.Writer.(*httptest.ResponseRecorder).Body
	var queryResponse client.CreateFlipsideQuerySuccessResponse
	if err = json.Unmarshal(bodyVerify.Bytes(), &queryResponse); err != nil {
		t.Fatalf("could not unmarshal response: %v", err)
	}

	paramTest := &client.GetFlipsideQueryResultRequest{
		Token: queryResponse.Token,
	}

	fmt.Println("found token ::: ", queryResponse.Token)

	// rest 60 seconds for waiting query result
	time.Sleep(60 * time.Second)

	err_token, ectx_token, hd_token := setupTest(t, http.MethodGet, cpkg.GetFlipsideQueryResultEndpoint, nil, &paramTest.Token)
	if err_token != nil {
		t.Fatalf("could not create handler: %v", err)
	}

	err_token = hd_token.GetFlipsideQueryResult(ectx_token)
	if err != nil {
		t.Fatalf("could not get query result: %v", err)
	}

	resToken := ectx_token.Response()
	if resToken.Status != http.StatusOK {
		t.Fatalf("expected status OK but got %v", resToken.Status)
	}

	bodyTokenVerify := resToken.Writer.(*httptest.ResponseRecorder).Body
	var queryTokenResponse client.GetFlipsideQueryResultSuccessResponse
	if err = json.Unmarshal(bodyTokenVerify.Bytes(), &queryTokenResponse); err != nil {
		t.Fatalf("could not unmarshal response: %v", err)
	}

	fmt.Println("queryTokenResponse", queryTokenResponse)
}

func TestEndToEnd_2_ReversedPython(t *testing.T) {
	bodyInPrompt := client.GPTPromptRequest{
		Prompt: "Find the transaction hash with the largest amount of ETH sent in the last 7 days.",
	}
	bodyInPromptRaw, err := json.Marshal(bodyInPrompt)
	if err != nil {
		t.Errorf("could not marshal request body: %v", err)
	}
	err_gpt, ectx_gpt, hd_gpt := setupTest(t, http.MethodGet, cpkg.GPTGenerateQueryEndpoint, &bodyInPromptRaw, nil)
	if err_gpt != nil {
		t.Fatalf("could not create handler: %v", err_gpt)
	}
	errGpt2 := hd_gpt.RunGptPythonClient(ectx_gpt)
	if errGpt2 != nil {
		t.Fatalf("could not create query: %v", errGpt2)
	}
	res_gpt := ectx_gpt.Response()
	if res_gpt.Status != http.StatusOK {
		t.Fatalf("expected status OK but got %v", res_gpt.Status)
	}
	fmt.Println("res_gpt", res_gpt)

	bodyverifyGpt := res_gpt.Writer.(*httptest.ResponseRecorder).Body
	var gptQueryResponse client.GPTPromptSuccessfulResponse
	if err_gpt = json.Unmarshal(bodyverifyGpt.Bytes(), &gptQueryResponse); err_gpt != nil {
		t.Fatalf("could not unmarshal response: %v", err_gpt)
	}

	bodyTest := &client.CreateFlipsideQueryRequest{
		Id:         string(rune(gptQueryResponse.Id)),
		Sql:        gptQueryResponse.Result,
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

	err, ectx, hd := setupTest(t, http.MethodGet, client.GetAllModels, &bodyRaw, nil)
	if err != nil {
		t.Fatalf("could not create handler: %v", err)
	}

	err = hd.CreateFlipsideQuery(ectx)
	if err != nil {
		t.Fatalf("could not create query: %v", err)
	}

	res := ectx.Response()
	if res.Status != http.StatusOK {
		t.Fatalf("expected status OK but got %v", res.Status)
	}

	bodyVerify := res.Writer.(*httptest.ResponseRecorder).Body
	var queryResponse client.CreateFlipsideQuerySuccessResponse
	if err = json.Unmarshal(bodyVerify.Bytes(), &queryResponse); err != nil {
		t.Fatalf("could not unmarshal response: %v", err)
	}

	paramTest := &client.GetFlipsideQueryResultRequest{
		Token: queryResponse.Token,
	}

	fmt.Println("found token ::: ", queryResponse.Token)

	// rest 60 seconds for waiting query result
	time.Sleep(60 * time.Second)

	err_token, ectx_token, hd_token := setupTest(t, http.MethodGet, cpkg.GetFlipsideQueryResultEndpoint, nil, &paramTest.Token)
	if err_token != nil {
		t.Fatalf("could not create handler: %v", err)
	}

	err_token = hd_token.GetFlipsideQueryResult(ectx_token)
	if err != nil {
		t.Fatalf("could not get query result: %v", err)
	}

	resToken := ectx_token.Response()
	if resToken.Status != http.StatusOK {
		t.Fatalf("expected status OK but got %v", resToken.Status)
	}

	bodyTokenVerify := resToken.Writer.(*httptest.ResponseRecorder).Body
	var queryTokenResponse client.GetFlipsideQueryResultSuccessResponse
	if err = json.Unmarshal(bodyTokenVerify.Bytes(), &queryTokenResponse); err != nil {
		t.Fatalf("could not unmarshal response: %v", err)
	}

	fmt.Println("queryTokenResponse", queryTokenResponse)
}

// INTEGRATION TEST USING OFFICIAL API FOR GPT 3.5 MODEL
func TestIntegration_1_OfficialClient(t *testing.T) {
	// STAGE 1. GPT
	//schemaRaw, err := engine.CreatePrompt("Find the most expensive gas fees-cost transaction in last 7 days.")
	bodyTest, err := client.NewChatCompletionRequest("Find the most expensive gas fees-cost transaction in last 7 days.", 3000, nil, nil, nil)
	if err != nil {
		t.Errorf("could not marshal request body: %v", err)
	}
	bodyRaw, err := json.Marshal(bodyTest)
	if err != nil {
		t.Errorf("could not marshal request body: %v", err)
	}
	err_gpt, ectx_gpt, hd_gpt := setupTest(t, http.MethodGet, cpkg.GPTGenerateQueryEndpoint, &bodyRaw, nil)
	if err_gpt != nil {
		t.Fatalf("could not create handler: %v", err_gpt)
	}
	errGpt2 := hd_gpt.CreateChatCompletion(ectx_gpt)
	if errGpt2 != nil {
		t.Fatalf("could not create query: %v", errGpt2)
	}
	res_gpt := ectx_gpt.Response()
	if res_gpt.Status != http.StatusOK {
		t.Fatalf("expected status OK but got %v", res_gpt.Status)
	}
	bodyverifyGpt := res_gpt.Writer.(*httptest.ResponseRecorder).Body
	var gptQueryResponse client.ChatCompletionResponse
	if err_gpt = json.Unmarshal(bodyverifyGpt.Bytes(), &gptQueryResponse); err_gpt != nil {
		t.Fatalf("could not unmarshal response: %v", err_gpt)
	}

	bodyTestCfqp := client.NewCreateFlipsideQueryResult(gptQueryResponse.Id, gptQueryResponse.GetContent())
	bodyRawCfqp, err := json.Marshal(bodyTestCfqp)
	if err != nil {
		t.Errorf("could not marshal request body: %v", err)
	}

	err, ectx, hd := setupTest(t, http.MethodGet, cpkg.GetFlipsideQueryResultEndpoint, &bodyRawCfqp, nil)
	if err != nil {
		t.Fatalf("could not create handler: %v", err)
	}
	err = hd.CreateFlipsideQuery(ectx)
	if err != nil {
		t.Fatalf("could not create query: %v", err)
	}
	res := ectx.Response()
	if res.Status != http.StatusOK {
		t.Fatalf("expected status OK but got %v", res.Status)
	}
	bodyVerify := res.Writer.(*httptest.ResponseRecorder).Body
	var queryResponse client.CreateFlipsideQuerySuccessResponse
	if err = json.Unmarshal(bodyVerify.Bytes(), &queryResponse); err != nil {
		t.Fatalf("could not unmarshal response: %v", err)
	}

	fmt.Println("queryResponse", queryResponse)

	paramTest := &client.GetFlipsideQueryResultRequest{
		Token: queryResponse.Token,
	}

	fmt.Println("found token ::: ", queryResponse.Token)

	// rest 60 seconds for waiting query result
	time.Sleep(60 * time.Second)

	err_token, ectx_token, hd_token := setupTest(t, http.MethodGet, cpkg.GetFlipsideQueryResultEndpoint, nil, &paramTest.Token)
	if err_token != nil {
		t.Fatalf("could not create handler: %v", err)
	}

	err_token = hd_token.GetFlipsideQueryResult(ectx_token)
	if err != nil {
		t.Fatalf("could not get query result: %v", err)
	}

	resToken := ectx_token.Response()
	if resToken.Status != http.StatusOK {
		t.Fatalf("expected status OK but got %v", resToken.Status)
	}

	bodyTokenVerify := resToken.Writer.(*httptest.ResponseRecorder).Body
	var queryTokenResponse client.GetFlipsideQueryResultSuccessResponse
	if err = json.Unmarshal(bodyTokenVerify.Bytes(), &queryTokenResponse); err != nil {
		t.Fatalf("could not unmarshal response: %v", err)
	}

	fmt.Println("queryTokenResponse", queryTokenResponse)
}
