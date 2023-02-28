package config

import (
	"fmt"
	"github.com/knadh/koanf"
)

type FlipsideEnv struct {
	API_KEY string `koanf:"%s.FlipsideEnv.API_KEY" envDefault:""`
}

func (fsenv *FlipsideEnv) ParseEnv(k *koanf.Koanf, env string) {
	fsenv.API_KEY = k.String(fmt.Sprintf("%s.FlipsideEnv.API_KEY", env))
}
