package models

import (
	"crow/oraiplayground/utils"
)

type AiSettings struct {
	Model             string
	Template          string
	MaxTokens         int
	Temperature       float64
	TopP              float64
	TopK              int
	FrequencyPenalty  float64
	PresencePenalty   float64
	RepetitionPenalty float64
}

func (s *AiSettings) SetParameters(params orai.Parameters) {
	s.MaxTokens = params.MaxTokens
	s.Temperature = params.Temperature
	s.TopP = params.TopP
	s.TopK = params.TopK
	s.FrequencyPenalty = params.FrequencyPenalty
	s.PresencePenalty = params.PresencePenalty
	s.RepetitionPenalty = params.RepetitionPenalty
}

func (s *AiSettings) GetParameters() orai.Parameters {
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

func (s *AiSettings) ParseFormData(form url.Values) (lastError error) {
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
