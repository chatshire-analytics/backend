package config

import (
	"fmt"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/labstack/gommon/log"
)

func LoadConfig(filePath string, env string) (*GlobalConfig, error) {
	k := koanf.New(".")

	if err := k.Load(file.Provider(filePath), yaml.Parser()); err != nil {
		log.Printf("could not load config file: %v", err)
		return nil, fmt.Errorf("error loading config file: %v", err)
	}

	var cfg GlobalConfig
	if err := k.Unmarshal(fmt.Sprintf("%s", env), &cfg); err != nil {
		log.Printf("could not unmarshal config file: %v", err)
		return nil, fmt.Errorf("error unmarshaling config file: %v", err)
	}
	cfg.OpenAIEnv.ParseEnv(k, env)

	return &cfg, nil
}
