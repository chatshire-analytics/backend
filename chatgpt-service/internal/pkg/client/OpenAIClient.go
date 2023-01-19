package client

import (
	"chatgpt-service/internal/pkg/engine"
	"net/http"
	"time"
)

type GPTClient struct {
	baseURL       string
	apiKey        string
	userAgent     string
	httpClient    *http.Client
	defaultEngine string
	idOrg         string
}

func NewGPTClient(apiKey string, options ...ClientOption) GPTClientInterface {
	httpClient := &http.Client{
		Timeout: time.Duration(engine.DefaultTimeoutSeconds * time.Second),
	}

	c := &GPTClient{
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
