package api

import (
	"chatgpt-service/internal/pkg/client"
	cif "chatgpt-service/pkg/client"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type Handler struct {
	oc   *client.OpenAIClient
	ectx *echo.Context
}

func NewHandler(c echo.Context) (*Handler, error) {
	ocInterface := c.Get(client.OpenAIClientKey)
	oc, ok := ocInterface.(*client.OpenAIClient)
	if !ok {
		return nil, errors.New("could not convert to OpenAI client")
	}
	return &Handler{
		oc:   oc,
		ectx: &c,
	}, nil
}

func (hd *Handler) ListModels(_ echo.Context) error {
	res, err := hd.oc.ListModels((*hd.ectx).Request().Context())
	if err != nil {
		return err
	}
	return (*hd.ectx).JSON(200, res)
}

func (hd *Handler) RetrieveModel(_ echo.Context) error {
	res, err := hd.oc.RetrieveModel((*hd.ectx).Request().Context(), (*hd.ectx).Param(cif.ModelIdParamKey))
	if err != nil {
		return err
	}
	return (*hd.ectx).JSON(200, res)
}

func (hd *Handler) CreateCompletion(_ echo.Context) error {
	var cr cif.CompletionRequest
	if err := (*hd.ectx).Bind(&cr); err != nil {
		return err
	}
	res, err := hd.oc.CreateCompletion((*hd.ectx).Request().Context(), cr)
	if err != nil {
		return err
	}
	return (*hd.ectx).JSON(200, res)
}

func (hd *Handler) CreateCompletionStream(_ echo.Context) error {
	var cr cif.CompletionRequest
	if err := (*hd.ectx).Bind(&cr); err != nil {
		(*hd.ectx).Error(err)
		return err
	}

	// Set up SSE
	(*hd.ectx).Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	(*hd.ectx).Response().Header().Set(echo.HeaderCacheControl, "no-cache")
	(*hd.ectx).Response().Header().Set(echo.HeaderConnection, "keep-alive")

	// Create a channel to receive new responses from the CompletionStream function
	respCh := make(chan cif.CompletionResponse)
	// A goroutine is started to run the CompletionStream function and send the received responses to the channel.
	go func() {
		err := hd.oc.CompletionStream((*hd.ectx).Request().Context(), cr, func(resp *cif.CompletionResponse) {
			respCh <- *resp
		})
		if err != nil {
			(*hd.ectx).Error(err)
		}
		close(respCh)
	}()

	// In the for-loop, the code continuously reads from the response channel
	// and sends updates to the client via SSE by writing to the response and flushing it
	// Continuously read from the response channel and send updates to the client via SSE
	_, err := (*hd.ectx).Response().Write([]byte("event: start"))
	if err != nil {
		(*hd.ectx).Error(err)
		return err
	}
	for {
		select {
		case resp, ok := <-respCh:
			if !ok {
				// Channel closed, done streaming
				_, err := (*hd.ectx).Response().Write([]byte("event: end"))
				if err != nil {
					return err
				}
				(*hd.ectx).Response().Flush()
				return nil
			}
			// Use SSE to stream updates to the client
			write, err := (*hd.ectx).Response().Write([]byte("data: " + resp.Choices[0].Text + "\n"))
			if err != nil {
				return err
			}
			(*hd.ectx).Response().Flush()
			if write == 0 {
				return nil
			}
		case <-(*hd.ectx).Request().Context().Done():
			// Request cancelled, done streaming
			write, err := (*hd.ectx).Response().Write([]byte("event: end"))
			if err != nil {
				return err
			}
			if write == 0 {
				return nil
			}
			(*hd.ectx).Response().Flush()
			return nil
		}
	}
}
