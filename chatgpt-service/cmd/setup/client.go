package setup

import (
	"chatgpt-service/internal/config"
	"chatgpt-service/internal/pkg/client"
	"chatgpt-service/internal/pkg/constants"
	"net/http"
	"time"
)

func NewOpenAIClient(cfg *config.GlobalConfig, options ...client.ClientOption) (client.OpenAIClientInterface, error) {
	httpClient := &http.Client{
		Timeout: time.Duration(constants.DefaultTimeoutSeconds * time.Second),
	}
	cl := &client.OpenAIClient{
		UserAgent:     constants.DefaultUserAgent,
		ApiKey:        cfg.OpenAIEnv.API_KEY,
		BaseURL:       constants.DefaultBaseURL,
		HttpClient:    httpClient,
		DefaultEngine: constants.DefaultEngine,
		IdOrg:         constants.DefaultUserName,
	}
	for _, clientOption := range options {
		err := clientOption(cl)
		if err != nil {
			return nil, err
		}
	}
	return cl, nil
}
