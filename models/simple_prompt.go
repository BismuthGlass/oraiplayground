package models

import (
	"crow/orai"
)

type SimplePrompt struct {
	Cue         string
	Context     string
	Instruction string
}

func (p *SimplePrompt) IntoLlama3Prompt() string {
	var messages []orai.ChatMessage
	if len(p.Instruction) > 0 {
		messages = append(messages, orai.ChatMessage{Role: orai.ChatRoleSystem, Message: p.Instruction})
	}
	if len(p.Context) > 0 {
		messages = append(messages, orai.ChatMessage{Role: orai.ChatRoleUser, Message: p.Context})
	}
	return orai.GenerateLlama3Prompt(messages, orai.ChatRoleAssistant, p.Cue)
}

func (p *SimplePrompt) IntoAlpacaPrompt() string {
	prompt := orai.AlpacaPrompt{
		Instruction: p.Instruction,
		Context:     p.Context,
		Cue:         p.Cue,
	}
	return prompt.String()
}
