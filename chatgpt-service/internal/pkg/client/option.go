package client

import (
	"net/http"
	"time"
)

// ClientOption are options that can be passed when creating a new client
type ClientOption func(*OpenAIClient) error

// WithOrg is a client option that allows you to override the organization ID
func WithOrg(id string) ClientOption {
	return func(c *OpenAIClient) error {
		c.IdOrg = id
		return nil
	}
}

// WithDefaultEngine is a client option that allows you to override the default client of the client
func WithDefaultEngine(engine string) ClientOption {
	return func(c *OpenAIClient) error {
		c.DefaultEngine = engine
		return nil
	}
}

// WithUserAgent is a client option that allows you to override the default user agent of the client
func WithUserAgent(userAgent string) ClientOption {
	return func(c *OpenAIClient) error {
		c.UserAgent = userAgent
		return nil
	}
}

// WithBaseURL is a client option that allows you to override the default base url of the client.
// The default base url is "https://api.openai.com/v1"
func WithBaseURL(baseURL string) ClientOption {
	return func(c *OpenAIClient) error {
		c.BaseURL = baseURL
		return nil
	}
}

// WithHTTPClient allows you to override the internal http.Client used
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *OpenAIClient) error {
		c.HttpClient = httpClient
		return nil
	}
}

// WithTimeout is a client option that allows you to override the default timeout duration of requests
// for the client. The default is 30 seconds. If you are overriding the http client as well, just include
// the timeout there.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *OpenAIClient) error {
		c.HttpClient.Timeout = timeout
		return nil
	}
}
