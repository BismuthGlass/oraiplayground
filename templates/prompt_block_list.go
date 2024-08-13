package templates

import (
	"crow/oraiplayground/models"
	"io"
	"log"
	"slices"
)

type promptBlockListItem struct {
	Name         string
	Active       bool
}

type promptBlockList struct {
	StoryName string
	Items []promptBlockListItem
}

func newPromptBlockList(story *models.Story) promptBlockList {
	preset := story.ActivePreset()
	items := make([]promptBlockListItem, 0, len(story.PromptBlocks))
	for _, b := range story.PromptBlocks {
		items = append(items, promptBlockListItem{
			Name: b.Name,
			Active: slices.Contains(preset.EnabledBlocks, b.Name),
		})
	}
	return promptBlockList{
		StoryName: story.Name,
		Items: items,
	}
}

func PromptBlockList(w io.Writer, story *models.Story) error {
	ctx := newPromptBlockList(story)
	err := engine.Template.ExecuteTemplate(w, "components/prompt_block_list.html", ctx)
	if err != nil {
		log.Println(err)
	}
	return err
}
