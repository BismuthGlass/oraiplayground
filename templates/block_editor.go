package templates

import (
	"crow/oraiplayground/models"
	"crow/oraiplayground/utils"
	"slices"
	"io"
	"log"
)

type BlockEditorListItem struct {
	Name string
	Favorite bool
	Enabled bool
}

type BlockEditorList struct {
	StoryName string
	Items []BlockEditorListItem
}

type BlockEditorForm struct {
	StoryName string
	Name string
	RoleOptions []utils.SelectOption
	Text string
	Message string
}

func NewBlockEditorList(story *models.Story) BlockEditorList {
	var items []BlockEditorListItem
	preset := story.ActivePreset()
	for _, b := range story.PromptBlocks {
		favorite := slices.Contains(preset.FavBlocks, b.Name)
		enabled := slices.Contains(preset.EnabledBlocks, b.Name)
		items = append(items, BlockEditorListItem{
			Name: b.Name,
			Favorite: favorite,
			Enabled: enabled,
		})
	}
	return BlockEditorList{
		StoryName: story.Name,
		Items: items,
	}
}

func NewBlockEditorForm(storyName string, block *models.PromptBlock, message string) BlockEditorForm {
	if block != nil {
		return BlockEditorForm{
			StoryName: storyName,
			Name: block.Name,
			RoleOptions: block.RoleOptions(),
			Text: block.Text,
			Message: message,
		}
	} else {
		return BlockEditorForm{
			StoryName: storyName,
			RoleOptions: models.PromptBlockRoleOptions(),
			Message: message,
		}
	}
}

func (e *Engine) BlockEditorList(w io.Writer, story *models.Story) error {
	ctx := NewBlockEditorList(story)
	err := e.Template.ExecuteTemplate(w, "components/block_editor_list.html", &ctx)
	if err != nil {
		log.Println(err)
	}
	return err
}

func (e *Engine) BlockEditorForm(w io.Writer, storyName string, block *models.PromptBlock, message string) error {
	ctx := NewBlockEditorForm(storyName, block, message)
	err := e.Template.ExecuteTemplate(w, "components/block_editor_form.html", &ctx)
	if err != nil {
		log.Println(err)
	}
	return err
}

