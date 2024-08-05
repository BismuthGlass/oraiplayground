package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"log"
	"crow/oraiplayground/utils"
)

type Config struct {
	OpenRouter struct {
		ApiKey string `yaml:"apiKey"`
	} `yaml:"openRouter"`
}

var ApiKey string
var AvailableModels []utils.SelectOption
var AvailableTemplates []utils.SelectOption

func LoadConfig() {
	data, err := os.ReadFile("config.yml")
	if err != nil {
		log.Fatal(err)
	}
	var config Config
	yaml.Unmarshal(data, &config)

	ApiKey = config.OpenRouter.ApiKey

	AvailableModels = []utils.SelectOption{
		{ Value: "lizpreciatior/lzlv-70b-fp16-hf" },
		{ Value: "meta-llama/llama-3-70b-instruct" },
		{ Value: "google/gemma-2-9b-it" },
		{ Value: "meta-llama/llama-3.1-70b-instruct" },
		{ Value: "meta-llama/llama-3.1-405b-instruct" },
	}

	AvailableTemplates = []utils.SelectOption{
		{ Value: "none", Name: "None" },
		{ Value: "alpaca", Name: "Alpaca" },
		{ Value: "llama3", Name: "Llama 3" },
		{ Value: "llama3_1", Name: "Llama 3.1" },
		{ Value: "gemma", Name: "Gemma" },
	}
}
