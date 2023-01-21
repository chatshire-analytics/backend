package api

import (
	"chatgpt-service/internal/pkg/client"
	cif "chatgpt-service/pkg/client"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

// TODO: refactor the code using a constructor function
//type Handler struct {
//	oc *client.OpenAIClient
//}
//
//func NewHandler(oc *client.OpenAIClient) *Handler {
//	return &Handler{
//		oc: oc,
//	}
//}

func ListModels(c echo.Context) error {
	ocInterface := c.Get(client.OpenAIClientKey)
	oc, ok := ocInterface.(*client.OpenAIClient)
	if !ok {
		return errors.New("could not convert to OpenAI client")
	}
	res, err := oc.ListModels(c.Request().Context())
	if err != nil {
		return err
	}
	return c.JSON(200, res)
}

func RetrieveModel(c echo.Context) error {
	ocInterface := c.Get(client.OpenAIClientKey)
	oc, ok := ocInterface.(*client.OpenAIClient)
	if !ok {
		return errors.New("could not convert to OpenAI client")
	}
	res, err := oc.RetrieveModel(c.Request().Context(), c.Param("model_id"))
	if err != nil {
		return err
	}
	return c.JSON(200, res)
}

func CreateCompletion(c echo.Context) error {
	ocInterface := c.Get(client.OpenAIClientKey)
	oc, ok := ocInterface.(*client.OpenAIClient)
	if !ok {
		return errors.New("could not convert to OpenAI client")
	}
	var cr cif.CompletionRequest
	if err := c.Bind(&cr); err != nil {
		return err
	}
	res, err := oc.CreateCompletion(c.Request().Context(), cr)
	if err != nil {
		return err
	}
	return c.JSON(200, res)
}

func CreateCompletionStream(c echo.Context) error {
	ocInterface := c.Get(client.OpenAIClientKey)
	oc, ok := ocInterface.(*client.OpenAIClient)
	if !ok {
		err := errors.New("could not convert to OpenAI client")
		c.Error(err)
		return err
	}

	var cr cif.CompletionRequest
	if err := c.Bind(&cr); err != nil {
		c.Error(err)
		return err
	}

	// Set up SSE
	c.Response().Header().Set(echo.HeaderContentType, "text/event-stream")
	c.Response().Header().Set(echo.HeaderCacheControl, "no-cache")
	c.Response().Header().Set(echo.HeaderConnection, "keep-alive")

	// Create a channel to receive new responses from the CompletionStream function
	respCh := make(chan cif.CompletionResponse)
	// A goroutine is started to run the CompletionStream function and send the received responses to the channel.
	go func() {
		err := oc.CompletionStream(c.Request().Context(), cr, func(resp *cif.CompletionResponse) {
			respCh <- *resp
		})
		if err != nil {
			c.Error(err)
		}
		close(respCh)
	}()

	// In the for-loop, the code continuously reads from the response channel
	// and sends updates to the client via SSE by writing to the response and flushing it
	// Continuously read from the response channel and send updates to the client via SSE
	_, err := c.Response().Write([]byte("data: \\start\n"))
	if err != nil {
		c.Error(err)
		return err
	}
	for {
		select {
		case resp, ok := <-respCh:
			if !ok {
				// Channel closed, done streaming
				_, err := c.Response().Write([]byte("data: \\end "))
				if err != nil {
					return err
				}
				c.Response().Flush()
				return nil
			}
			// Use SSE to stream updates to the client
			write, err := c.Response().Write([]byte("data: " + resp.Choices[0].Text + "\n"))
			if err != nil {
				return err
			}
			c.Response().Flush()
			if write == 0 {
				return nil
			}
		case <-c.Request().Context().Done():
			// Request cancelled, done streaming
			write, err := c.Response().Write([]byte("data: \\end"))
			if err != nil {
				return err
			}
			if write == 0 {
				return nil
			}
			c.Response().Flush()
			return nil
		}
	}

	/*
		print each stream onto the console
			onData := func(resp *cif.CompletionResponse) {
				fmt.Println(resp.Choices[0].Text)
			}
			err := oc.CompletionStream(c.Request().Context(), cr, onData)
			if err != nil {
				c.Error(err)
				return err
			}
			return c.JSON(http.StatusOK, "OK")
	*/
}
