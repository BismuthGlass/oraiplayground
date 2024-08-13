package templates

import (
	"crow/oraiplayground/models"
	"io"
	"log"
)

type PromptInfo struct {
	StoryName string
	LastPromptTokens int
	LastPrompt string
	NextPromptTokens int
	NextPrompt string
}

func (e *Engine) PromptInfo(w io.Writer, story *models.Story) error {
	ctx := PromptInfo {
		StoryName: story.Name,
		LastPrompt: "This was the last prompt you issued!",
		NextPrompt: "This is the next prompt you're going to issue!",
	}
	err := e.Template.ExecuteTemplate(w, "components/prompt_info.html", ctx)
	if err != nil {
		log.Println(err)
	}
	return err
}

