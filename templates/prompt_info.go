package templates

import (
	"crow/oraiplayground/models"
	"crow/oraiplayground/utils"
	"io"
	"log"
	"math"
)

type promptInfo struct {
	StoryName string
	LastPromptTokens int
	LastPrompt string
	NextPromptTokens int
	NextPrompt string
}

func PromptInfo(w io.Writer, story *models.Story) error {
	nextPrompt := story.GenPrompt()
	nextPromptTokens := int(math.Ceil(float64(utils.WordCount(nextPrompt)) * (1 / 0.75)))
	ctx := promptInfo {
		StoryName: story.Name,
		LastPrompt: story.LastPrompt,
		LastPromptTokens: story.LastPromptTokens,
		NextPrompt: nextPrompt,
		NextPromptTokens: nextPromptTokens,
	}
	err := engine.Template.ExecuteTemplate(w, "components/prompt_info.html", ctx)
	if err != nil {
		log.Println(err)
	}
	return err
}

