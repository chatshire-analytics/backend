package api

import (
	"chatgpt-service/internal/config"
	"chatgpt-service/internal/pkg/client"
	"chatgpt-service/internal/pkg/store"
	cif "chatgpt-service/pkg/client"
	"fmt"
	"github.com/labstack/echo/v4"
	"html/template"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

type Handler struct {
	oc *client.OpenAIClient
	fc *client.FlipsideClient
	// TODO: remove echo.Context to have different context for each request
	ectx *echo.Context
	// TODO: remove handler - database mapping connection
	db *store.Database
}

func NewHandler(c echo.Context, cfg config.GlobalConfig, oc *client.OpenAIClient, fc *client.FlipsideClient, db *store.Database) (*Handler, error) {
	return &Handler{
		oc:   oc,
		ectx: &c,
		fc:   fc,
		db:   db,
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

func (hd *Handler) TempRequestNewQuery(c echo.Context) error {
	var request cif.GPTPromptRequest
	if err := c.Bind(&request); err != nil {
		return c.JSON(400, err.Error())
	}
	cr, err := cif.NewChatCompletionRequest(request.Prompt, 3000, nil, nil, nil)
	if err != nil {
		return (*hd.ectx).JSON(400, err.Error())
	}
	crResult, err := hd.oc.CreateNewChatCompletion(c.Request().Context(), *cr)
	if err != nil {
		return (*hd.ectx).JSON(503, err.Error())
	}
	err = hd.db.StoreGptSqlResult(*cr, crResult)
	if err != nil {
		return (*hd.ectx).JSON(500, err.Error())
	}
	cq := cif.NewCreateFlipsideQueryResult(crResult.Id, crResult.GetContent())
	if err != nil {
		return (*hd.ectx).JSON(500, err.Error())
	}
	res, err := hd.fc.CreateFlipsideQuery(c.Request().Context(), *cq)
	resBody := make(map[string]string)
	if err != nil {
		// TODO: temporarily return the error message as the response body
		resBody["status"] = err.Error()
		return (*hd.ectx).JSON(http.StatusBadRequest, resBody)
	}
	err = hd.db.UpdateCreateFlipsideQueryResult(cq.Id, res.Token)
	if err != nil {
		resBody["status"] = err.Error()
		fmt.Println(err.Error())
		return (*hd.ectx).JSON(http.StatusServiceUnavailable, resBody)
	}

	// rest for 5 seconds
	time.Sleep(5 * time.Second)

	fqr := cif.NewGetFlipsideQueryResultRequest(res.Token)
	result, err := hd.fc.GetFlipsideQueryResult(c.Request().Context(), *fqr)
	ulResBody := make(map[string]interface{})
	if err != nil {
		if strings.Contains(err.Error(), "running") {
			ulResBody["token"] = res.Token
			ulResBody["status"] = "running! if it takes too long, submit new query"
			return c.JSON(http.StatusAccepted, ulResBody)
		}
		return err
	}
	err = hd.db.StoreGetFlipsideQueryResult(*fqr, *result)
	if err != nil || result == nil {
		resBody["status"] = err.Error()
		fmt.Println(err.Error())
		return c.JSON(http.StatusServiceUnavailable, resBody)
	}
	ulResBody["results"] = result.Results
	ulResBody["query"] = crResult.GetContent()
	ulResBody["status"] = "success"

	return c.JSON(http.StatusOK, ulResBody)
}

func (hd *Handler) TempHTMLHandler(c echo.Context) error {
	data := struct {
		Title string
	}{
		Title: "Chatshire Beta - On-Chain Data Analytics Tool Using GPT Model",
	}

	tmpl, err := template.ParseFiles("./template/index.html")
	if err != nil {
		return c.JSON(500, err.Error())
	}

	return tmpl.Execute(c.Response(), data)
}

func (hd *Handler) CreateChatCompletion(_ echo.Context) error {
	var cr cif.ChatCompletionRequest
	if err := (*hd.ectx).Bind(&cr); err != nil {
		return (*hd.ectx).JSON(400, err.Error())
	}
	res, err := hd.oc.CreateNewChatCompletion((*hd.ectx).Request().Context(), cr)
	if err != nil {
		return (*hd.ectx).JSON(503, err.Error())
	}
	err = hd.db.StoreGptSqlResult(cr, res)
	if err != nil {
		return (*hd.ectx).JSON(500, err.Error())
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

func (hd *Handler) RunGptPythonClient(_ echo.Context) error {
	accessToken, err := (*hd.oc).GetAccessToken()
	if err != nil {
		return err
	}

	var promptRaw cif.GPTPromptRequest
	if err := (*hd.ectx).Bind(&promptRaw); err != nil {
		return (*hd.ectx).JSON(400, err.Error())
	}
	// TODO: temporarily
	prompt, err := cif.CreatePrompt(promptRaw)
	if err != nil {
		return err
	}
	promptInString := prompt.String()

	result, err := exec.Command("python", "../pkg/client/ChatbotRunner.py", accessToken, promptInString).Output()
	if err != nil {
		return (*hd.ectx).JSON(500, err.Error())
	}

	id, err := hd.db.StoreGptPythonSqlResult(promptRaw.Prompt, string(result))
	if err != nil {
		return err
	}

	responseBody := cif.GPTPromptSuccessfulResponse{
		Id:     id,
		Result: string(result),
	}

	return (*hd.ectx).JSON(200, responseBody)
}

func (hd *Handler) CreateFlipsideQuery(_ echo.Context) error {
	var cq cif.CreateFlipsideQueryRequest
	if err := (*hd.ectx).Bind(&cq); err != nil {
		return err
	}
	res, err := hd.fc.CreateFlipsideQuery((*hd.ectx).Request().Context(), cq)
	resBody := make(map[string]string)
	if err != nil {
		// TODO: temporarily return the error message as the response body
		resBody["status"] = err.Error()
		return (*hd.ectx).JSON(http.StatusBadRequest, resBody)
	}
	err = hd.db.UpdateCreateFlipsideQueryResult(cq.Id, res.Token)
	if err != nil {
		resBody["status"] = err.Error()
		fmt.Println(err.Error())
		return (*hd.ectx).JSON(http.StatusServiceUnavailable, resBody)
	}
	return (*hd.ectx).JSON(200, res)
}

func (hd *Handler) GetFlipsideQueryResult(ctx echo.Context) error {
	token := ctx.Param("token")
	gr := cif.GetFlipsideQueryResultRequest{
		Token: token,
	}
	result, err := hd.fc.GetFlipsideQueryResult((*hd.ectx).Request().Context(), gr)
	// TODO: extract the response body separately. the response body should be {"status": "running! if it takes too long, submit new query", "token": token}
	resBody := make(map[string]string)
	if err != nil {
		if strings.Contains(err.Error(), "running") {
			resBody["token"] = token
			resBody["status"] = "running! if it takes too long, submit new query"
			return (*hd.ectx).JSON(http.StatusAccepted, resBody)
		}
		return err
	}
	err = hd.db.StoreGetFlipsideQueryResult(gr, *result)
	if err != nil {
		resBody["status"] = err.Error()
		fmt.Println(err.Error())
		return (*hd.ectx).JSON(http.StatusServiceUnavailable, resBody)
	}
	return (*hd.ectx).JSON(http.StatusOK, result)
}
