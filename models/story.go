package models

import (
	"errors"
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
		{ Name: "block_1", Text: "" },
		{ Name: "block_2", Text: "" },
		{ Name: "block_3", Text: "" },
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
	// TODO
	return nil
}

func (s *Story) AddPromptBlock(block PromptBlock) error {
	existingBlock := s.GetPromptBlock(block.Name)
	if existingBlock != nil {
		return errors.New("Exists")
	}
	s.PromptBlocks = append(s.PromptBlocks, block)
	return nil
}

func (s *Story) DeletePromptBlock(name string) {
	// TODO
}

// Moves `name` under `ref`
func (s *Story) MovePromptBlock(name string, ref string) error {
	// TODO
	return nil
}

func (s *Story) GetEnabledBlocks() []PromptBlock {
	//preset := s.ActivePreset()
	return nil
}

