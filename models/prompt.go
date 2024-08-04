package models

import (
	"bytes"
	"crow/orai"
	"crow/oraiplayground/utils"
	"errors"
	"fmt"
	"slices"
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
	Section  PromptSection  `json:"section"`
	Role     PromptRole     `json:"role"`
	IsGroup  bool           `json:"isGroup"`
	Active   bool           `json:"active"`
	Text     string         `json:"text"`
	Prefix   string         `json:"prefix"`
	Suffix   string         `json:"suffix"`
	Children []*PromptBlock `json:"children"`
}

func (pb *PromptBlock) SectionOptions() []utils.SelectOption {
	options := []utils.SelectOption{
		{ Value: PromptSectionInstruction, Name: "Instruction" },
		{ Value: PromptSectionContext, Name: "Context" },
		{ Value: PromptSectionCue, Name: "Cue" },
	}
	utils.SetSelection(options, string(pb.Section))
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
	utils.SetSelection(options, string(pb.Section))
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

func (root *PromptBlock) getBlockLocation(name string) (array *[]*PromptBlock, index int) {
	for i, b := range root.Children {
		if b.Name == name {
			return &root.Children, i
		}
		if b.IsGroup {
			array, index := b.getBlockLocation(name)
			if array != nil {
				return array, index
			}
		}
	}

	return nil, -1
}

func (root *PromptBlock) blockRef(name string) *PromptBlock {
	a, i := root.getBlockLocation(name)
	if a == nil {
		return nil
	}
	return (*a)[i]
}

func (root *PromptBlock) GetBlock(name string) (PromptBlock, bool) {
	r := root.blockRef(name)
	if r == nil {
		return PromptBlock{}, false
	}
	return *r, true
}

func (root *PromptBlock) TraverseTree(callback func(*PromptBlock) error) error {
	for _, b := range root.Children {
		err := callback(b)
		if err != nil {
			return err
		}
		if b.IsGroup {
			err = b.TraverseTree(callback)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (root *PromptBlock) ChangeLocation(target string, destinationGroup string, destinationReference string) error {
	targetLocation, targetIndex := root.getBlockLocation(target)
	if targetLocation == nil {
		return errors.New("target not found")
	}
	targetBlock := (*targetLocation)[targetIndex]
	(*targetLocation) = slices.Delete(*targetLocation, targetIndex, targetIndex + 1)

	if len(destinationGroup) > 0 && len(destinationReference) == 0 {
		newParent := root.blockRef(destinationGroup)
		if newParent == nil {
			return errors.New("destination group not found")
		}
		newParent.Children = append(newParent.Children, targetBlock)
		return nil
	}

	referenceLocation, referenceIndex := root.getBlockLocation(destinationReference)
	if referenceLocation == nil {
		return errors.New("reference block not found")
	}
	(*referenceLocation) = slices.Insert(*referenceLocation, referenceIndex, targetBlock)
	return nil
}

func (root *PromptBlock) UpdateBlock(new PromptBlock) error {
	target := root.blockRef(new.Name)
	if target == nil {
		return errors.New("target not found")
	}

	target.Text = new.Text
	target.Prefix = new.Prefix
	target.Suffix = new.Suffix
	target.Section = new.Section
	return nil
}

func (root *PromptBlock) SetActive(name string, active bool) error {
	target := root.blockRef(name)
	if target == nil {
		return errors.New("target not found")
	}

	target.Active = active
	return nil
}

// Returns all the removed blocks
func (root *PromptBlock) RemoveBlocks(names []string) []*PromptBlock {
	var removed []*PromptBlock
	for _, n := range names {
		loc, i := root.getBlockLocation(n)
		if loc == nil {
			continue
		}
		target := (*loc)[i]
		(*loc) = append((*loc)[:i], (*loc)[i+1:]...)

		removed = append(removed, target.Children...)
		removed = append(removed, target)
	}
	return removed
}

func (root *PromptBlock) PromptText() string {
	if !root.Active {
		return ""
	}
	if !root.IsGroup {
		return root.Text
	}

	var text bytes.Buffer
	text.WriteString(root.Prefix)
	for _, n := range root.Children {
		text.WriteString(n.PromptText())
	}
	text.WriteString(root.Suffix)
	return text.String()
}

func (root *PromptBlock) IntoSimplePrompt() SimplePrompt {
	var instruction bytes.Buffer
	var context bytes.Buffer
	var cue bytes.Buffer

	for _, b := range root.Children {
		if !b.Active {
			continue
		}
		text := b.PromptText()
		switch b.Section {
		case "instruction":
			instruction.WriteString(text)
		case "context":
			context.WriteString(text)
		case "cue":
			cue.WriteString(text)
		default:
			fmt.Println("Warning, unknown section:", b.Section)
		}
	}

	return SimplePrompt{
		Instruction: instruction.String(),
		Context: context.String(),
		Cue: cue.String(),
	}
}

func (root *PromptBlock) ActivateBlock(name string, active bool) error {
	block := root.blockRef(name)
	if block == nil {
		return errors.New("not found")
	}
	block.Active = active
	return nil
}

func (root *PromptBlock) GatherSection(section PromptSection) string {
	var b bytes.Buffer
	for _, c := range root.Children {
		if c.Section == section {
			b.WriteString(c.PromptText())
		}
	}
	return b.String()
}

func (root *PromptBlock) GatherRole(role PromptRole) string {
	var b bytes.Buffer
	for _, c := range root.Children {
		if c.Role == role {
			b.WriteString(c.PromptText())
		}
	}
	return b.String()
}

func (root *PromptBlock) IntoPrompt(template string) string {
	switch template {
	case "alpaca":
		return root.IntoAlpacaPrompt()
	case "llama3":
		return root.IntoLlama3Prompt()
	case "llama3_1":
		return root.IntoLlama3_1Prompt()
	case "gemma":
		return root.IntoGemmaPrompt()
	default:
		return root.IntoRawPrompt()
	}
}

func (root *PromptBlock) IntoMessages() []orai.ChatMessage {
	var messages []orai.ChatMessage
	for _, b := range root.Children {
		switch b.Role {
		case PromptRoleUser:
			messages = append(messages, orai.ChatMessage{
				Role: orai.ChatRoleUser,
				Message: b.PromptText(),
			})
		case PromptRoleAssistant:
			messages = append(messages, orai.ChatMessage{
				Role: orai.ChatRoleAssistant,
				Message: b.PromptText(),
			})
		case PromptRoleSystem:
			messages = append(messages, orai.ChatMessage{
				Role: orai.ChatRoleSystem,
				Message: b.PromptText(),
			})
		}
	}
	return messages
}

func (root *PromptBlock) IntoAlpacaPrompt() string {
	simplePrompt := root.IntoSimplePrompt()
	return simplePrompt.IntoAlpacaPrompt()
}

func (root *PromptBlock) IntoLlama3Prompt() string {
	simplePrompt := root.IntoSimplePrompt()
	return simplePrompt.IntoLlama3Prompt()
}

func (root *PromptBlock) IntoLlama3_1Prompt() string {
	simplePrompt := root.IntoSimplePrompt()
	return simplePrompt.IntoLlama3_1Prompt()
}

func (root *PromptBlock) IntoGemmaPrompt() string {
	var turns []orai.ChatMessage

	for _, block := range root.Children {
		if !block.Active {
			continue
		}
		var role orai.ChatRole
		switch block.Section {
		case "instruction", "context":
			role = orai.ChatRoleUser
		default:
			role = orai.ChatRoleAssistant
		}
		turns = append(turns, orai.ChatMessage{
			Role: role,
			Message: block.PromptText(),
		})
	}

	turns = orai.CollapseMessages(turns)
	return orai.GenerateGemmaPrompt(turns, orai.ChatRoleAssistant, "")
}

func (root *PromptBlock) IntoRawPrompt() string {
	var buffer bytes.Buffer
	for _, block := range root.Children {
		buffer.WriteString(block.PromptText())
	}
	return buffer.String()
}
