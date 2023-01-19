package client

import (
	"chatgpt-service/internal/pkg/engine"
	"chatgpt-service/pkg/client"
	"context"
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

// TODO: implement above referrencing https://github.com/PullRequestInc/go-gpt3/blob/main/gpt3.go
func (G OpenAIClient) ListModels(ctx context.Context) (*client.ListModelsResponse, error) {
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
		idOrg:         "",
	}
	for _, o := range options {
		err := o(c)
		if err != nil {
			return nil
		}
	}
	return c
}
