package models

import (
	"errors"
	"slices"
)

type StoryMode string
const (
	StoryModeInstruct = "instruct"
	StoryModeChat = "chat"
)

type Story struct {
	Name string
	Description string
	Mode StoryMode
	ModelSettings ModelSettings

	LastPrompt string
	LastPromptTokens int

	ActivePromptPreset string
	PromptPresets map[string]*PromptSettings
	PromptBlocks []PromptBlock
	Variables []StoryVariable
}

func (s *Story) compileBlocks() []PromptBlock {
	return s.PromptBlocks
}

func (s *Story) GetPromptBlock(name string) *PromptBlock {
	for i, _ := range s.PromptBlocks {
		if s.PromptBlocks[i].Name == name {
			return &s.PromptBlocks[i]
		}
	}
	return nil
}

func (s *Story) GenPrompt() string {
	blocks := s.compileBlocks()
	return BlocksIntoPrompt(blocks, s.ModelSettings.Template)
}

func NewStory(name string, description string, storyMode StoryMode, defaultModelSettings ModelSettings) Story {
	blocks := []PromptBlock {
		{ Name: "block_1", Text: "", Role: PromptRoleUser },
		{ Name: "block_2", Text: "", Role: PromptRoleUser },
		{ Name: "block_3", Text: "", Role: PromptRoleUser },
	}
	promptPresets := make(map[string]*PromptSettings)
	promptPresets["default"] = &PromptSettings{}
	return Story{
		Name: name,
		Description: description,
		Mode: storyMode,
		ModelSettings: defaultModelSettings,
		ActivePromptPreset: "default",
		PromptBlocks: blocks,
		PromptPresets: promptPresets,
	}
}

func (s *Story) ActivePreset() *PromptSettings {
	return s.PromptPresets[s.ActivePromptPreset]
}

func (s *Story) UpdatePreset(name string, data PromptSettings) {
	s.PromptPresets[name] = &data
}

func (s *Story) EnablePromptBlock(name string) error {
	block := s.GetPromptBlock(name)
	if block == nil {
		return errors.New("not found")
	}
	s.ActivePreset().enableBlock(name)
	return nil
}

func (s *Story) DisablePromptBlock(name string) error {
	block := s.GetPromptBlock(name)
	if block == nil {
		return errors.New("not found")
	}
	s.ActivePreset().disableBlock(name)
	return nil
}

func (s *Story) TogglePromptBlockFavorite(name string) error {
	block := s.GetPromptBlock(name)
	if block == nil {
		return errors.New("not found")
	}
	s.ActivePreset().toggleBlockFavorite(name)
	return nil
}

func (s *Story) TogglePromptBlock(name string) error {
	block := s.GetPromptBlock(name)
	if block == nil {
		return errors.New("not found")
	}
	s.ActivePreset().toggleBlock(name)
	return nil
}

func (s *Story) UpdatePromptBlock(name string, blockData PromptBlock) error {
	block := s.GetPromptBlock(name)
	if block == nil {
		return errors.New("not found")
	}
	existingBlock := s.GetPromptBlock(blockData.Name)
	if existingBlock != nil && existingBlock != block {
		return errors.New("exists")
	}
	block.Name = blockData.Name
	block.Role = blockData.Role
	block.Text = blockData.Text
	block.Compiled = blockData.Compiled
	return nil
}

func (s *Story) AddPromptBlock(block PromptBlock) error {
	existingBlock := s.GetPromptBlock(block.Name)
	if existingBlock != nil {
		return errors.New("exists")
	}
	s.PromptBlocks = append(s.PromptBlocks, block)
	return nil
}

func (s *Story) DeletePromptBlock(name string) {
	s.PromptBlocks = slices.DeleteFunc(s.PromptBlocks, func(b PromptBlock) bool {
		return b.Name == name
	})
	for _, preset := range s.PromptPresets {
		preset.removeBlock(name)
	}
}

// Moves `name` under `ref`
func (s *Story) MovePromptBlock(name string, ref string) error {
	blockPtr := s.GetPromptBlock(name)
	if blockPtr == nil || s.GetPromptBlock(ref) == nil {
		return errors.New("not found")
	}
	if name == ref {
		return nil
	}
	block := *blockPtr

	s.PromptBlocks = slices.DeleteFunc(s.PromptBlocks, func(b PromptBlock) bool { return b.Name == name })
	for i, b := range s.PromptBlocks {
		if b.Name == ref {
			s.PromptBlocks = slices.Insert(s.PromptBlocks, i + 1, block)
			return nil
		}
	}
	return nil
}

func (s *Story) GetEnabledBlocks() []PromptBlock {
	preset := s.ActivePreset()
	list := make([]PromptBlock, 0, len(preset.FavBlocks))
	for _, b := range s.PromptBlocks {
		if slices.Contains(preset.EnabledBlocks, b.Name) {
			list = append(list, b)
		}
	}
	return list
}

