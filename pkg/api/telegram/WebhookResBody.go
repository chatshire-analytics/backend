package telegram

type WebhookResBody struct {
	ChatId int64  `json:"chat_id"`
	Text   string `json:"text"`
}
