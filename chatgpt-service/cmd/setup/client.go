package setup

import (
	"chatgpt-service/internal/config"
	"chatgpt-service/internal/pkg/client"
	"chatgpt-service/internal/pkg/engine"
	"net/http"
	"time"
)

func NewOpenAIClient(cfg *config.GlobalConfig, options ...client.ClientOption) (client.OpenAIClientInterface, error) {
	httpClient := &http.Client{
		Timeout: time.Duration(engine.DefaultTimeoutSeconds * time.Second),
	}
	cl := &client.OpenAIClient{
		UserAgent:     engine.DefaultUserAgent,
		ApiKey:        cfg.OpenAIEnv.API_KEY,
		BaseURL:       engine.DefaultBaseURL,
		HttpClient:    httpClient,
		DefaultEngine: engine.DefaultEngine,
		IdOrg:         engine.DefaultUserName,
	}
	for _, clientOption := range options {
		err := clientOption(cl)
		if err != nil {
			return nil, err
		}
	}
	return cl, nil
}
