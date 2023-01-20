package client

import (
	"bytes"
	"chatgpt-service/pkg/client"
	cerror "chatgpt-service/pkg/errors"
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
)

const OpenAIClientKey = "OpenAIClient"

type OpenAIClient struct {
	BaseURL       string
	ApiKey        string
	UserAgent     string
	HttpClient    *http.Client
	DefaultEngine string
	IdOrg         string
}

func (oc *OpenAIClient) JSONBodyReader(body interface{}) (io.Reader, error) {
	if body == nil {
		return bytes.NewBuffer(nil), nil
	}
	raw, err := json.Marshal(body)
	if err != nil {
		return nil, errors.New("failed to encode body: " + err.Error())
	}
	return bytes.NewBuffer(raw), nil
}

func (oc *OpenAIClient) NewRequestBuilder(ctx context.Context, method string, path string, payload interface{}) (*http.Request, error) {
	br, err := oc.JSONBodyReader(payload)
	if err != nil {
		return nil, err
	}
	url := oc.BaseURL + path // link to openai.com
	req, err := http.NewRequestWithContext(ctx, method, url, br)
	if err != nil {
		return nil, err
	}
	if len(oc.IdOrg) > 0 {
		req.Header.Set("user", oc.IdOrg)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", oc.ApiKey))
	return req, nil
}

func (oc *OpenAIClient) ExecuteRequest(req *http.Request) (*http.Response, error) {
	resp, err := oc.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if err := oc.CheckRequestSucceed(resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (oc *OpenAIClient) CheckRequestSucceed(resp *http.Response) error {
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

func (oc *OpenAIClient) getResponseObject(rsp *http.Response, v interface{}) error {
	defer rsp.Body.Close()
	if err := json.NewDecoder(rsp.Body).Decode(v); err != nil {
		return fmt.Errorf("invalid json response: %w", err)
	}
	return nil
}

func (oc *OpenAIClient) ListModels(ctx context.Context) (*client.ListModelsResponse, error) {
	endPoint := client.ModelEndPoint
	req, err := oc.NewRequestBuilder(ctx, http.MethodGet, endPoint, nil)
	if err != nil {
		return nil, err
	}
	resp, err := oc.ExecuteRequest(req)
	if err != nil {
		return nil, err
	}

	output := new(client.ListModelsResponse)
	if err := oc.getResponseObject(resp, output); err != nil {
		return nil, err
	}
	return output, nil
}

func (oc OpenAIClient) RetrieveModel(ctx context.Context, engine string) (*client.ModelObject, error) {
	req, err := oc.NewRequestBuilder(ctx, http.MethodGet, client.ModelEndPoint+"/"+engine, nil)
	if err != nil {
		return nil, err
	}
	resp, err := oc.ExecuteRequest(req)
	if err != nil {
		return nil, err
	}
	output := new(client.ModelObject)
	if err := oc.getResponseObject(resp, output); err != nil {
		return nil, err
	}
	return output, nil
}

func (oc OpenAIClient) Completion(ctx context.Context, request client.CompletionRequest) (*client.CompletionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (oc OpenAIClient) CompletionStream(ctx context.Context, request client.CompletionRequest, onData func(response *client.CompletionResponse)) error {
	//TODO implement me
	panic("implement me")
}

func (oc OpenAIClient) CompletionWithEngine(ctx context.Context, engine string, request client.CompletionRequest) (*client.CompletionResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (oc OpenAIClient) CompletionStreamWithEngine(ctx context.Context, engine string, request client.CompletionRequest, onData func(response *client.CompletionResponse)) error {
	//TODO implement me
	panic("implement me")
}

func (oc OpenAIClient) Edits(ctx context.Context, request client.EditsRequest) (*client.EditsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (oc OpenAIClient) Embeddings(ctx context.Context, request client.EmbeddingsRequest) (*client.EmbeddingsResponse, error) {
	//TODO implement me
	panic("implement me")
}
