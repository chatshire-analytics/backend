package main

import (
	"mentat-backend/gateway-service/cmd/setup"
)

func main() {
	err := setup.InitializeEcho().Start(":8080")
	if err != nil {
		return
	}

	err = setup.SetupTelegramBot()
	if err != nil {
		return
	}
}
