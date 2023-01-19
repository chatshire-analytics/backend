package client

import (
	"net/http"
	"time"
)

type GPTClientType struct {
	baseURL       string
	apiKey        string
	userAgent     string
	httpClient    *http.Client
	defaultEngine string
	idOrg         string
}

func NewGPTClient(apiKey string, options ...ClientOption) GPTClientInterface {
	httpClient := &http.Client{
		Timeout: time.Duration(defaultTimeoutSeconds * time.Second),
	}

	c := &GPTClientType{
		userAgent:     defaultUserAgent,
		apiKey:        apiKey,
		baseURL:       defaultBaseURL,
		httpClient:    httpClient,
		defaultEngine: DefaultEngine,
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
