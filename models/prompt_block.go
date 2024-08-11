package models

import (
	"bytes"
	"crow/orai"
	"crow/oraiplayground/utils"
)

type PromptSection string
const (
	PromptSectionInstruction = "instruction"
	PromptSectionContext = "context"
	PromptSectionCue = "cue"
)

type PromptRole string
const (
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
}

func (pb *PromptBlock) SectionOptions() []utils.SelectOption {
	options := []utils.SelectOption{
		{ Value: PromptSectionInstruction, Name: "Instruction" },
		{ Value: PromptSectionContext, Name: "Context" },
		{ Value: PromptSectionCue, Name: "Cue" },
	}
	//utils.SetSelection(options, string(pb.Section))
	return options
}

func (pb *PromptBlock) RoleOptions() []utils.SelectOption {
	options := []utils.SelectOption{
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
		{ Value: orai.ChatRoleUser, Name: "User" },
		{ Value: orai.ChatRoleAssistant, Name: "Assistant" },
		{ Value: orai.ChatRoleSystem, Name: "System" },
	}
	utils.SetSelection(options, string(orai.ChatRoleUser))
	return options
}

func PromptBlockSectionOptions() []utils.SelectOption {
	options := []utils.SelectOption{
		{ Value: PromptSectionInstruction, Name: "Instruction" },
		{ Value: PromptSectionContext, Name: "Context" },
		{ Value: PromptSectionCue, Name: "Cue" },
	}
	utils.SetSelection(options, string(PromptSectionInstruction))
	return options
}

func BlocksIntoSimplePrompt(blocks []PromptBlock) SimplePrompt {
	var instruction bytes.Buffer
	var context bytes.Buffer
	var cue bytes.Buffer

	//for _, b := range blocks {
	//	text := b.Text
	//	switch b.Section {
	//	case "instruction":
	//		instruction.WriteString(text)
	//	case "context":
	//		context.WriteString(text)
	//	case "cue":
	//		cue.WriteString(text)
	//	default:
	//		fmt.Println("Warning, unknown section:", b.Section)
	//	}
	//}

	return SimplePrompt{
		Instruction: instruction.String(),
		Context: context.String(),
		Cue: cue.String(),
	}
}

func PromptGatherSection(blocks []PromptBlock, section PromptSection) string {
	var b bytes.Buffer
	//for _, c := range blocks {
	//	if c.Section == section {
	//		b.WriteString(c.Text)
	//	}
	//}
	return b.String()
}

func PromptGatherRole(blocks []PromptBlock, role PromptRole) string {
	var b bytes.Buffer
	for _, c := range blocks {
		if c.Role == role {
			b.WriteString(c.Text)
		}
	}
	return b.String()
}

func BlocksIntoPrompt(blocks []PromptBlock, template string) string {
	switch template {
	case "alpaca":
		return BlocksIntoAlpacaPrompt(blocks)
	case "llama3":
		return BlocksIntoLlama3Prompt(blocks)
	case "llama3_1":
		return BlocksIntoLlama3_1Prompt(blocks)
	case "gemma":
		return BlocksIntoGemmaPrompt(blocks)
	default:
		return BlocksIntoRawPrompt(blocks)
	}
}

func BlocksIntoMessages(blocks []PromptBlock) []orai.ChatMessage {
	var messages []orai.ChatMessage
	//for _, b := range blocks {
	//	switch b.Role {
	//	case PromptRoleUser:
	//		messages = append(messages, orai.ChatMessage{
	//			Role: orai.ChatRoleUser,
	//			Message: b.Text,
	//		})
	//	case PromptRoleAssistant:
	//		messages = append(messages, orai.ChatMessage{
	//			Role: orai.ChatRoleAssistant,
	//			Message: b.Text,
	//		})
	//	case PromptRoleSystem:
	//		messages = append(messages, orai.ChatMessage{
	//			Role: orai.ChatRoleSystem,
	//			Message: b.Text,
	//		})
	//	}
	//}
	return messages
}

func BlocksIntoAlpacaPrompt(blocks []PromptBlock) string {
	simplePrompt := BlocksIntoSimplePrompt(blocks)
	return simplePrompt.IntoAlpacaPrompt()
}

func BlocksIntoLlama3Prompt(blocks []PromptBlock) string {
	simplePrompt := BlocksIntoSimplePrompt(blocks)
	return simplePrompt.IntoLlama3Prompt()
}

func BlocksIntoLlama3_1Prompt(blocks []PromptBlock) string {
	simplePrompt := BlocksIntoSimplePrompt(blocks)
	return simplePrompt.IntoLlama3_1Prompt()
}

func BlocksIntoGemmaPrompt(blocks []PromptBlock) string {
	var turns []orai.ChatMessage

	//for _, block := range blocks {
	//	var role orai.ChatRole
	//	switch block.Section {
	//	case "instruction", "context":
	//		role = orai.ChatRoleUser
	//	default:
	//		role = orai.ChatRoleAssistant
	//	}
	//	turns = append(turns, orai.ChatMessage{
	//		Role: role,
	//		Message: block.Text,
	//	})
	//}

	turns = orai.CollapseMessages(turns)
	return orai.GenerateGemmaPrompt(turns, orai.ChatRoleAssistant, "")
}

func BlocksIntoRawPrompt(blocks []PromptBlock) string {
	var buffer bytes.Buffer
	for _, block := range blocks {
		buffer.WriteString(block.Text)
	}
	return buffer.String()
}
