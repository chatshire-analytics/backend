package client

import (
	"bytes"
	"chatgpt-service/internal/pkg/engine"
	"chatgpt-service/pkg/client"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"time"
)

type OpenAIClient struct {
	baseURL       string
	apiKey        string
	userAgent     string
	httpClient    *http.Client
	defaultEngine string
	idOrg         string
}

func (gc *OpenAIClient) JSONBodyReader(body interface{}) (io.Reader, error) {
	if body == nil {
		return bytes.NewBuffer(nil), nil
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return nil, errors.New("failed to encode body: " + err.Error())
	}
	return bytes.NewBuffer(raw), nil
}

func (gc *OpenAIClient) NewRequest(ctx context.Context, method string, path string, payload interface{}) (*http.Request, error) {
	br, err := gc.JSONBodyReader(payload)
	if err != nil {
		return nil, err
	}
	url := gc.baseURL + path // link to openai.com
	req, err := http.NewRequestWithContext(ctx, method, url, br)
	if err != nil {
		return nil, err
	}
	if len(gc.idOrg) > 0 {
		req.Header.Set("user", gc.idOrg)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", gc.apiKey))
	return req, nil
}

// TODO: implement above referrencing https://github.com/PullRequestInc/go-gpt3/blob/main/gpt3.go
func (gc *OpenAIClient) ListModels(ctx context.Context) (*client.ListModelsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (G OpenAIClient) RetrieveModel(ctx context.Context, engine string) (*client.ModelObject, error) {
	//TODO implement me
	panic("implement me")
}

func (G OpenAIClient) Completion(ctx context.Context, request client.CompletionRequest) (*client.CompletionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (G OpenAIClient) CompletionStream(ctx context.Context, request client.CompletionRequest, onData func(response *client.CompletionResponse)) error {
	//TODO implement me
	panic("implement me")
}

func (G OpenAIClient) CompletionWithEngine(ctx context.Context, engine string, request client.CompletionRequest) (*client.CompletionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (G OpenAIClient) CompletionStreamWithEngine(ctx context.Context, engine string, request client.CompletionRequest, onData func(response *client.CompletionResponse)) error {
	//TODO implement me
	panic("implement me")
}

func (G OpenAIClient) Edits(ctx context.Context, request client.EditsRequest) (*client.EditsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (G OpenAIClient) Embeddings(ctx context.Context, request client.EmbeddingsRequest) (*client.EmbeddingsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func NewGPTClient(apiKey string, options ...ClientOption) GPTClientInterface {
	httpClient := &http.Client{
		Timeout: time.Duration(engine.DefaultTimeoutSeconds * time.Second),
	}

	c := &OpenAIClient{
		userAgent:     engine.DefaultUserAgent,
		apiKey:        apiKey,
		baseURL:       engine.DefaultBaseURL,
		httpClient:    httpClient,
		defaultEngine: engine.DefaultEngine,
		idOrg:         engine.DefaultUserName,
	}
	for _, o := range options {
		err := o(c)
		if err != nil {
			return nil
		}
	}
	return c
}
