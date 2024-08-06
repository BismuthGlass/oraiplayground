package templates

import (
	"crow/oraiplayground/models"
	"slices"
	"io"
	"log"
)

type PromptBlockEditorRow struct {
	Name string
	Favorite bool
	Enabled bool
}

type PromptBlockEditorTable struct {
	StoryName string
	Rows []PromptBlockEditorRow
}

type PromptBlockEditor struct {
	Mode string
	Table PromptBlockEditorTable
	Row PromptBlockEditorRow
}

func NewPromptBlockEditorTable(story *models.Story) PromptBlockEditor {
	var rows []PromptBlockEditorRow
	preset := story.ActivePreset()
	for _, b := range story.PromptBlocks {
		favorite := slices.Contains(preset.FavBlocks, b.Name)
		enabled := slices.Contains(preset.EnabledBlocks, b.Name)
		rows = append(rows, PromptBlockEditorRow{
			Name: b.Name,
			Favorite: favorite,
			Enabled: enabled,
		})
	}
	return PromptBlockEditor{
		Mode: "table",
		Table: PromptBlockEditorTable{
			StoryName: story.Name,
			Rows: rows,
		},
	}
}

func (e *Engine) PromptBlockEditorTable(w io.Writer, story *models.Story) error {
	ctx := NewPromptBlockEditorTable(story)
	err := e.Template.ExecuteTemplate(w, "components/prompt_block_editor.html", &ctx)
	if err != nil {
		log.Println(err)
	}
	return err
}
