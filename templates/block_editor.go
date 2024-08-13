package templates

import (
	"crow/oraiplayground/models"
	"crow/oraiplayground/utils"
	"slices"
	"io"
	"log"
)

type blockEditorListItem struct {
	Name string
	Favorite bool
	Enabled bool
}

type blockEditorList struct {
	StoryName string
	Items []blockEditorListItem
}

type blockEditorForm struct {
	StoryName string
	Name string
	RoleOptions []utils.SelectOption
	Text string
	Message string
}

func newBlockEditorList(story *models.Story) blockEditorList {
	var items []blockEditorListItem
	preset := story.ActivePreset()
	for _, b := range story.PromptBlocks {
		favorite := slices.Contains(preset.FavBlocks, b.Name)
		enabled := slices.Contains(preset.EnabledBlocks, b.Name)
		items = append(items, blockEditorListItem{
			Name: b.Name,
			Favorite: favorite,
			Enabled: enabled,
		})
	}
	return blockEditorList{
		StoryName: story.Name,
		Items: items,
	}
}

func newBlockEditorForm(storyName string, block *models.PromptBlock, message string) blockEditorForm {
	if block != nil {
		return blockEditorForm{
			StoryName: storyName,
			Name: block.Name,
			RoleOptions: block.RoleOptions(),
			Text: block.Text,
			Message: message,
		}
	} else {
		return blockEditorForm{
			StoryName: storyName,
			RoleOptions: models.PromptBlockRoleOptions(),
			Message: message,
		}
	}
}

func BlockEditorList(w io.Writer, story *models.Story) error {
	ctx := newBlockEditorList(story)
	err := engine.Template.ExecuteTemplate(w, "components/block_editor_list.html", &ctx)
	if err != nil {
		log.Println(err)
	}
	return err
}

func BlockEditorForm(w io.Writer, storyName string, block *models.PromptBlock, message string) error {
	ctx := newBlockEditorForm(storyName, block, message)
	err := engine.Template.ExecuteTemplate(w, "components/block_editor_form.html", &ctx)
	if err != nil {
		log.Println(err)
	}
	return err
}

