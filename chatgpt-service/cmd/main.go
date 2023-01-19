package main

import (
	"fmt"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

var k = koanf.New(".")

func main() {
	err := k.Load(file.Provider("config.yaml"), yaml.Parser())
	if err != nil {
		return
	}
	fmt.Println(k.String("dev.OpenAI.API_KEY"))
}
