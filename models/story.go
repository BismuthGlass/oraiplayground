package models

import (
	"crow/orai"
	"crow/oraiplayground/utils"
	"net/url"
)

type StoryMode string
const (
	StoryModeText = "text"
	StoryModeChat = "chat"
	StoryModeInstruct = "instruct"
)

type StorySettings struct {
	Model             string
	Template          string
	Mode              StoryMode
	MaxTokens         int
	Temperature       float64
	TopP              float64
	TopK              int
	FrequencyPenalty  float64
	PresencePenalty   float64
	RepetitionPenalty float64
}

func (s *StorySettings) SetParameters(params orai.Parameters) {
	s.MaxTokens = params.MaxTokens
	s.Temperature = params.Temperature
	s.TopP = params.TopP
	s.TopK = params.TopK
	s.FrequencyPenalty = params.FrequencyPenalty
	s.PresencePenalty = params.PresencePenalty
	s.RepetitionPenalty = params.RepetitionPenalty
}

func (s *StorySettings) GetParameters() orai.Parameters {
	return orai.Parameters{
		MaxTokens:         s.MaxTokens,
		Temperature:       s.Temperature,
		TopP:              s.TopP,
		TopK:              s.TopK,
		FrequencyPenalty:  s.FrequencyPenalty,
		PresencePenalty:   s.PresencePenalty,
		RepetitionPenalty: s.RepetitionPenalty,
	}
}

func (s *StorySettings) ParseFormData(form url.Values) (lastError error) {
	if form.Has("model") {
		s.Model = form.Get("model")
	}
	if form.Has("template") {
		s.Template = form.Get("template")
	}
	if form.Has("maxTokens") {
		value, err := utils.ParseAndValidateInt("max tokens", form.Get("maxTokens"), 0, 1000000)
		if err != nil {
			lastError = err
		} else {
			s.MaxTokens = int(value)
		}
	}
	if form.Has("temperature") {
		value, err := utils.ParseAndValidateFloat("temperature", form.Get("temperature"), 0.0, 2.0)
		if err != nil {
			lastError = err
		} else {
			s.Temperature = value
		}
	}
	if form.Has("topP") {
		value, err := utils.ParseAndValidateFloat("top P", form.Get("topP"), 0.0, 1.0)
		if err != nil {
			lastError = err
		} else {
			s.TopP = value
		}
	}
	if form.Has("topK") {
		value, err := utils.ParseAndValidateInt("top K", form.Get("topK"), 0, 1000000)
		if err != nil {
			lastError = err
		} else {
			s.TopK = int(value)
		}
	}
	if form.Has("frequencyPenalty") {
		value, err := utils.ParseAndValidateFloat("frequency penalty", form.Get("frequencyPenalty"), -2.0, 2.0)
		if err != nil {
			lastError = err
		} else {
			s.FrequencyPenalty = value
		}
	}
	if form.Has("presencePenalty") {
		value, err := utils.ParseAndValidateFloat("presence penalty", form.Get("presencePenalty"), -2.0, 2.0)
		if err != nil {
			lastError = err
		} else {
			s.PresencePenalty = value
		}
	}
	if form.Has("repetitionPenalty") {
		value, err := utils.ParseAndValidateFloat("repetition penalty", form.Get("repetitionPenalty"), 0.0, 2.0)
		if err != nil {
			lastError = err
		} else {
			s.RepetitionPenalty = value
		}
	}
	return
}

type Story struct {
	Id    int64
	Title string
	Settings StorySettings
	PromptRoot *PromptBlock
	ChatMessages []orai.ChatMessage
	Text string
}

func NewStory(id int64, title string) Story {
	settings := StorySettings{}
	settings.SetParameters(orai.DefaultAIParameters())
	return Story{
		Settings: settings,
		PromptRoot: &PromptBlock{},
	}
}