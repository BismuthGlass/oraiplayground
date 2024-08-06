package templates

import (
	"crow/oraiplayground/models"
	"crow/oraiplayground/utils"
	"crow/oraiplayground/config"
	"io"
	"log"
)

type StorySettings struct {
	Models            []utils.SelectOption
	Templates         []utils.SelectOption
	MaxTokens         int
	Temperature       float64
	TopP              float64
	TopK              int
	FrequencyPenalty  float64
	PresencePenalty   float64
	RepetitionPenalty float64
	LastError         error
}

func NewStorySettings(story *models.Story, err error) StorySettings {
	models := make([]utils.SelectOption, len(config.AvailableModels))
	_ = copy(models, config.AvailableModels)
	utils.SetSelection(models, story.ModelSettings.Model)

	aiTemplates := make([]utils.SelectOption, len(config.AvailableTemplates))
	_ = copy(aiTemplates, config.AvailableTemplates)
	utils.SetSelection(aiTemplates, story.ModelSettings.Template)
	
	return StorySettings{
		Models: models,
		Templates: aiTemplates,
		MaxTokens: story.ModelSettings.MaxTokens,
		Temperature: story.ModelSettings.Temperature,
		TopP: story.ModelSettings.TopP,
		TopK: story.ModelSettings.TopK,
		FrequencyPenalty: story.ModelSettings.FrequencyPenalty,
		PresencePenalty: story.ModelSettings.PresencePenalty,
		RepetitionPenalty: story.ModelSettings.RepetitionPenalty,
		LastError: err,
	}
}

func (e *Engine) Settings(w io.Writer, ctx *StorySettings) error {
 	err := e.Template.ExecuteTemplate(w, "components/settings.html", ctx)
	if err != nil {
		log.Println(err)
	}
	return err
}

type StoryLayout struct {
	Type string
}

type Story struct {
	Mode models.StoryMode
	Settings StorySettings
	PromptBlockList PromptBlockList
	PromptBlockEditor PromptBlockEditor
}

func NewStory(story *models.Story) Story {
	return Story{
		Mode: story.Mode,
		Settings: NewStorySettings(story, nil),
		PromptBlockList: NewPromptBlockList(story),
		PromptBlockEditor: NewPromptBlockEditorTable(story),
	}
}

func (e *Engine) StoryPage(w io.Writer, story *models.Story) error {
	ctx := NewStory(story)
	err := e.Template.ExecuteTemplate(w, "page_story.html", &ctx)
	if err != nil {
		log.Println(err)
	}
	return err
}
