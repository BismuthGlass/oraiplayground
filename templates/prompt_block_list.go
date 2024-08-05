package templates

import (
	"crow/oraiplayground/models"
	"io"
	"log"
	"slices"
)

type PromptBlockListItem struct {
	Name         string
	Active       bool
}

type PromptBlockList struct {
	StoryName string
	Items []PromptBlockListItem
}

func NewPromptBlockList(story *models.Story) PromptBlockList {
	preset := story.ActivePreset()
	items := make([]PromptBlockListItem, 0, len(story.PromptBlocks))
	for _, b := range story.PromptBlocks {
		items = append(items, PromptBlockListItem{
			Name: b.Name,
			Active: slices.Contains(preset.EnabledBlocks, b.Name),
		})
	}
	return PromptBlockList{
		StoryName: story.Name,
		Items: items,
	}
}

func (e *Engine) PromptBlockList(w io.Writer, ctx *PromptBlockList) error {
	err := e.Template.ExecuteTemplate(w, "components/prompt_block_list.html", ctx)
	if err != nil {
		log.Println(err)
	}
	return err
}
