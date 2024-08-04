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

func (p *SimplePrompt) IntoLlama3_1Prompt() string {
	var turns []orai.Llama3_1Turn
	if len(p.Instruction) > 0 {
		turns = append(turns, orai.Llama3_1Turn{Role: orai.Llama3_1RoleSystem, Message: p.Instruction})
	}
	if len(p.Context) > 0 {
		turns = append(turns, orai.Llama3_1Turn{Role: orai.Llama3_1RoleUser, Message: p.Context})
	}
	turns = append(turns, orai.Llama3_1Turn{Role: orai.Llama3_1RoleAssistant, Message: p.Cue})
	return orai.Llama3_1CompilePrompt(turns)
}

func (p *SimplePrompt) IntoAlpacaPrompt() string {
	prompt := orai.AlpacaPrompt{
		Instruction: p.Instruction,
		Context:     p.Context,
		Cue:         p.Cue,
	}
	return prompt.String()
}
