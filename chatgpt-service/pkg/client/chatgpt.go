package client

type GPTPromptRequest struct {
	Prompt string
}

type GPTPromptSuccessfulResponse struct {
	Id     int
	Result string
}
