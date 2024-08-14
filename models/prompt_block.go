package models

import (
	"bytes"
	"crow/orai/templates"
	"crow/oraiplayground/utils"
)

type PromptRole string
const (
	PromptRoleNone = "none"
	PromptRoleUser = "user"
	PromptRoleAssistant = "assistant"
	PromptRoleSystem = "system"
	PromptRoleUserCue = "userCue"
	PromptRoleAssistantCue = "assistantCue"
)

type PromptBlock struct {
	Name     string         `json:"name"`
	Role     PromptRole     `json:"role"`
	Text     string         `json:"text"`
	Compiled bool           `json:"compiled"`
}

func (pb *PromptBlock) RoleOptions() []utils.SelectOption {
	options := []utils.SelectOption{
		{ Value: PromptRoleNone, Name: "None" },
		{ Value: PromptRoleUser, Name: "User" },
		{ Value: PromptRoleAssistant, Name: "Assistant" },
		{ Value: PromptRoleSystem, Name: "System" },
		{ Value: PromptRoleAssistantCue, Name: "Assistant Cue" },
		{ Value: PromptRoleUserCue, Name: "User Cue" },
	}
	utils.SetSelection(options, string(pb.Role))
	return options
}

func PromptBlockRoleOptions() []utils.SelectOption {
	options := []utils.SelectOption{
		{ Value: PromptRoleNone, Name: "None" },
		{ Value: PromptRoleUser, Name: "User" },
		{ Value: PromptRoleAssistant, Name: "Assistant" },
		{ Value: PromptRoleSystem, Name: "System" },
		{ Value: PromptRoleAssistantCue, Name: "Assistant Cue" },
		{ Value: PromptRoleUserCue, Name: "User Cue" },
	}
	utils.SetSelection(options, string(PromptRoleUser))
	return options
}

func BlocksIntoRawPrompt(blocks []PromptBlock) string {
	var prompt bytes.Buffer
	for _, b := range blocks {
		prompt.WriteString(b.Text)
	}
	return prompt.String()
}

func BlocksIntoAlpaca(blocks []PromptBlock) string {
	var instruction bytes.Buffer
	var context bytes.Buffer
	var cue bytes.Buffer

	for _, b := range blocks {
		switch b.Role {
		case PromptRoleSystem:
			instruction.WriteString(b.Text)
		case PromptRoleUser, PromptRoleUserCue:
			context.WriteString(b.Text)
		case PromptRoleAssistant, PromptRoleAssistantCue:
			cue.WriteString(b.Text)
		default:
		}
	}

	prompt := templates.AlpacaPrompt{
		Instruction: instruction.String(),
		Context: context.String(),
		Cue: cue.String(),
	}
	return prompt.Prompt(true)
}

func BlocksIntoPrompt(blocks []PromptBlock, template string) string {
	switch template {
	case "alpaca":
		return BlocksIntoAlpaca(blocks)
	//case "llama3":
	//	return BlocksIntoLlama3Prompt(blocks)
	//case "llama3_1":
	//	return BlocksIntoLlama3_1Prompt(blocks)
	//case "gemma":
	//	return BlocksIntoGemmaPrompt(blocks)
	default:
		return BlocksIntoRawPrompt(blocks)
	}
}

//func BlocksIntoLlama3_1(blocks []PromptBlock) string {
//	var messages []orai.Llama3_1Turn
//}

//func BlocksIntoLlama3_1Prompt(blocks []PromptBlock) string {
//	simplePrompt := BlocksIntoSimplePrompt(blocks)
//	return simplePrompt.IntoLlama3_1Prompt()
//}

