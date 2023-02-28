package client

import (
	"bytes"
	cerror "chatgpt-service/pkg/errors"
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
)

const FlipsideClientKey = "FlipsideClient"

type FlipsideClient struct {
	BaseURL    string
	ApiKey     string
	UserAgent  string
	HttpClient *http.Client
}

func (fc *FlipsideClient) JSONBodyReader(body interface{}) (io.Reader, error) {
	if body == nil {
		return bytes.NewBuffer(nil), nil
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return nil, errors.New("failed to encode body: " + err.Error())
	}
	return bytes.NewBuffer(raw), nil
}

func (fc *FlipsideClient) NewRequestBuilder(ctx context.Context, method string, path string, payload interface{}) (*http.Request, error) {
	br, err := fc.JSONBodyReader(payload)
	if err != nil {
		return nil, err
	}
	url := fc.BaseURL + path // link to openai.com
	req, err := http.NewRequestWithContext(ctx, method, url, br)
	if err != nil {
		return nil, err
	}
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", fc.ApiKey))
	return req, nil
}

func (fc *FlipsideClient) ExecuteRequest(req *http.Request) (*http.Response, error) {
	resp, err := fc.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if err := fc.CheckRequestSucceed(resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (fc *FlipsideClient) CheckRequestSucceed(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read from body: %w", err)
	}
	var result cerror.APIErrorResponse
	if err := json.Unmarshal(data, &result); err != nil {
		// if we can't decode the json error then create an unexpected error
		apiError := cerror.APIError{
			StatusCode: resp.StatusCode,
			Type:       "Unexpected",
			Message:    string(data),
		}
		return apiError
	}
	result.Error.StatusCode = resp.StatusCode
	return result.Error
}

func (fc *FlipsideClient) getResponseObject(rsp *http.Response, v interface{}) error {
	defer rsp.Body.Close()
	if err := json.NewDecoder(rsp.Body).Decode(v); err != nil {
		return fmt.Errorf("invalid json response: %w", err)
	}
	return nil
}
