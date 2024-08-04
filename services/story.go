package services

import (
	"crow/orai"
	"crow/oraiplayground/components"
	"crow/oraiplayground/models"
	"crow/oraiplayground/templates"
	"crow/oraiplayground/utils"
	"errors"
)

const StoryCtxKey = utils.CtxKey("ServiceStory")

type Story struct {
	Story models.Story
	AiServer *AiServer
	PromptListComponent components.PromptList
	AvailableModels []utils.SelectOption
	AvailableTemplates []utils.SelectOption
}

func NewStory(aiServer *AiServer) *Story {
	blocks := []*models.PromptBlock{
		{
			Name: "instruction",
			Section: "instruction",
			Active: true,
		},
		{
			Name: "context",
			Section: "context",
			Active: true,
		},
		{
			Name: "cue",
			Section: "cue",
			Active: true,
		},
	}

	storyModel := models.Story{
		PromptRoot: &models.PromptBlock{
			Children: blocks,
		},
	}
	storyModel.Settings.SetParameters(orai.DefaultAIParameters())
	storyModel.Settings.Mode = models.StoryModeInstruct
	storyModel.Settings.Model = "lizpreciatior/lzlv-70b-fp16-hf"
	storyModel.Settings.Template = "alpaca"

	promptListComponent := components.PromptList{
		DataRoot: storyModel.PromptRoot,
	}

	return &Story{
		Story: storyModel,
		AiServer: aiServer,
		PromptListComponent: promptListComponent,
		AvailableModels: []utils.SelectOption{
			{ Value: "lizpreciatior/lzlv-70b-fp16-hf" },
			{ Value: "meta-llama/llama-3-70b-instruct" },
			{ Value: "google/gemma-2-9b-it" },
			{ Value: "meta-llama/llama-3.1-70b-instruct" },
			{ Value: "meta-llama/llama-3.1-405b-instruct" },
		},
		AvailableTemplates: []utils.SelectOption{
			{ Value: "none", Name: "None" },
			{ Value: "alpaca", Name: "Alpaca" },
			{ Value: "llama3", Name: "Llama 3" },
			{ Value: "llama3_1", Name: "Llama 3.1" },
			{ Value: "gemma", Name: "Gemma" },
		},
	}
}

func (s *Story) SettingsTmplModel(err error) templates.StorySettings {
	models := make([]utils.SelectOption, len(s.AvailableModels))
	_ = copy(models, s.AvailableModels)
	utils.SetSelection(models, s.Story.Settings.Model)

	aiTemplates := make([]utils.SelectOption, len(s.AvailableTemplates))
	_ = copy(aiTemplates, s.AvailableTemplates)
	utils.SetSelection(aiTemplates, s.Story.Settings.Template)
	
	return templates.StorySettings{
		Models: models,
		Templates: aiTemplates,
		MaxTokens: s.Story.Settings.MaxTokens,
		Temperature: s.Story.Settings.Temperature,
		TopP: s.Story.Settings.TopP,
		TopK: s.Story.Settings.TopK,
		FrequencyPenalty: s.Story.Settings.FrequencyPenalty,
		PresencePenalty: s.Story.Settings.PresencePenalty,
		RepetitionPenalty: s.Story.Settings.RepetitionPenalty,
		LastError: err,
	}
}

func (s *Story) NewPromptBlock(block models.PromptBlock) error {
	_, found := s.Story.PromptRoot.GetBlock(block.Name)
	if found {
		return errors.New("Exists")
	}
	s.Story.PromptRoot.Children = append(s.Story.PromptRoot.Children, &block)
	return nil
}

func (s *Story) DeletePromptBlocks(names []string) {
	removed := s.Story.PromptRoot.RemoveBlocks(names)
	for _, b := range removed {
		delete(s.PromptListComponent.Collapsed, b.Name)
	}
}

func (s *Story) GetPrompt() string {
	return s.Story.PromptRoot.IntoPrompt(s.Story.Settings.Template)
}

func (s *Story) RequestPrompt(continueCue string) AiServiceRequestClientInfo {
	return s.AiServer.IssueRequest(
		s.Story.Settings.Model, 
		s.Story.Settings.GetParameters(),
		s.Story.PromptRoot.IntoPrompt(s.Story.Settings.Template) + continueCue,
	)
}

func (s *Story) GetPromptRequest(id int64) AiServiceResponse {
	return <-s.AiServer.RequestChannel(id)
}

func (s *Story) CancelPromptRequest(id int64) error {
	return s.AiServer.CancelRequest(id)
}

func (s *Story) ActivatePromptBlock(block string, active bool) error {
	return s.Story.PromptRoot.ActivateBlock(block, active)
}

func (s *Story) AddMessage(role orai.ChatRole, message string) {
	s.Story.ChatMessages = append(s.Story.ChatMessages, orai.ChatMessage{
		Role: role,
		Message: message,
	})
}

func (s *Story) RequestMessage(role orai.ChatRole, message string) AiServiceRequestClientInfo {
	blockMessages := s.Story.PromptRoot.IntoMessages()
	var cue string
	switch role {
	case orai.ChatRoleUser:
		cue = s.Story.PromptRoot.GatherRole(models.PromptRoleUserCue) + message
	case orai.ChatRoleAssistant:
		cue = s.Story.PromptRoot.GatherRole(models.PromptRoleAssistantCue) + message
	default:
		cue = message
	}
	prompt := orai.GenerateLlama3Prompt(blockMessages, role, cue)
	return s.AiServer.IssueRequest(
		s.Story.Settings.Model, 
		s.Story.Settings.GetParameters(),
		prompt,
	)
}
