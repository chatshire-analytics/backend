package config

import (
	"fmt"
	"github.com/knadh/koanf"
)

type OpenAIENV struct {
	API_KEY string `koanf:"%s.OpenAIEnv.API_KEY" envDefault:""`
}

func (oaenv *OpenAIENV) ParseEnv(k *koanf.Koanf, env string) {
	oaenv.API_KEY = k.String(fmt.Sprintf("%s.OpenAIEnv.API_KEY", env))
}
