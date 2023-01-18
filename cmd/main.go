package main

import (
	"mentat-backend/cmd/setup"
)

func main() {
	err := setup.InitializeEcho().Start(":8080")
	if err != nil {
		return
	}
}
