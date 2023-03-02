package client

type GPTPromptRequest struct {
	Prompt string `json:"prompt"`
}

type GPTPromptSuccessfulResponse struct {
	Id     int    `json:"id"`
	Result string `json:"result"`
}
