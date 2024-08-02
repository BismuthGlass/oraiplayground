package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"log"
)

type Config struct {
	OpenRouter struct {
		ApiKey string `yaml:"apiKey"`
	} `yaml:"openRouter"`
}

var ApiKey string

func LoadConfig() {
	data, err := os.ReadFile("config.yml")
	if err != nil {
		log.Fatal(err)
	}
	var config Config
	yaml.Unmarshal(data, &config)

	ApiKey = config.OpenRouter.ApiKey
}
