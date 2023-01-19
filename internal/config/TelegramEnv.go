package config

import (
	"fmt"
	"github.com/knadh/koanf"
)

type TelegramEnv struct {
	API_ID  string `koanf:"%s.TelegramEnv.API_ID" envDefault:""`
	API_KEY string `koanf:"%s.TelegramEnv.API_KEY" envDefault:""`
}

func (tgenv *TelegramEnv) ParseEnv(k *koanf.Koanf, env string) {
	tgenv.API_ID = k.String(fmt.Sprintf("%s.TelegramEnv.API_ID", env))
	tgenv.API_KEY = k.String(fmt.Sprintf("%s.TelegramEnv.API_KEY", env))
}
