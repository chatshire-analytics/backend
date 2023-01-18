package setup

import "github.com/pkg/errors"

const TelegramErrorString = "unexpecte Http Status Code from Telegram API"

func TelegramError(httpStatus string) error {
	return errors.New(TelegramErrorString + " httpStatus: " + httpStatus)
}
