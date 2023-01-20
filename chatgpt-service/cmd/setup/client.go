package setup

import (
	"chatgpt-service/internal/pkg/client"
	"chatgpt-service/internal/pkg/engine"
	"net/http"
	"time"
)

func NewOpenAIClient(apiKey string, options ...client.ClientOption) client.OpenAIClientInterface {
	httpClient := &http.Client{
		Timeout: time.Duration(engine.DefaultTimeoutSeconds * time.Second),
	}

	cl := &client.OpenAIClient{
		UserAgent:     engine.DefaultUserAgent,
		ApiKey:        apiKey,
		BaseURL:       engine.DefaultBaseURL,
		HttpClient:    httpClient,
		DefaultEngine: engine.DefaultEngine,
		IdOrg:         engine.DefaultUserName,
	}
	for _, clientOption := range options {
		err := clientOption(cl)
		if err != nil {
			return nil
		}
	}
	return cl
}
