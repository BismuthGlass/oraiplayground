package templates

import (
	"crow/oraiplayground/models"
	"crow/oraiplayground/utils"
	"crow/oraiplayground/config"
	"io"
	"log"
)

type storySettings struct {
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

type story struct {
	StoryName string
	Mode models.StoryMode
	Settings storySettings
	PlaygroundBlockList blockEditorList
	BlockEditorForm blockEditorForm
	BlockEditorList blockEditorList
}

func newStorySettings(story *models.Story, err error) storySettings {
	models := make([]utils.SelectOption, len(config.AvailableModels))
	_ = copy(models, config.AvailableModels)
	utils.SetSelection(models, story.ModelSettings.Model)

	aiTemplates := make([]utils.SelectOption, len(config.AvailableTemplates))
	_ = copy(aiTemplates, config.AvailableTemplates)
	utils.SetSelection(aiTemplates, story.ModelSettings.Template)
	
	return storySettings{
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

func newStory(s *models.Story) story {
	return story{
		StoryName: s.Name,
		Mode: s.Mode,
		Settings: newStorySettings(s, nil),
		PlaygroundBlockList: newPlaygroundBlockList(s),
		BlockEditorForm: newBlockEditorForm(s.Name, nil, ""),
		BlockEditorList: newBlockEditorList(s),
	}
}

func StorySettings(w io.Writer, story *models.Story, errorInfo error) error {
	ctx := newStorySettings(story, errorInfo)
 	err := engine.Template.ExecuteTemplate(w, "components/settings.html", ctx)
	if err != nil {
		log.Println(err)
	}
	return err
}

func StoryPage(w io.Writer, story *models.Story) error {
	ctx := newStory(story)
	err := engine.Template.ExecuteTemplate(w, "page_story.html", &ctx)
	if err != nil {
		log.Println(err)
	}
	return err
}
