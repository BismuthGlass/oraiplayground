package templates

import (
	"crow/oraiplayground/models"
	"crow/oraiplayground/utils"
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
	PromptBlockEditors []PromptBlockEditor
}

func (e *Engine) StoryPage(w io.Writer, ctx *Story) error {
	err := e.Template.ExecuteTemplate(w, "page_story.html", ctx)
	if err != nil {
		log.Println(err)
	}
	return err
}
