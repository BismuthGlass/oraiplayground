package templates

import (
	"crow/oraiplayground/models"
	"crow/oraiplayground/utils"
	"io"
)

type PromptBlockEditor struct {
	Mode models.StoryMode
	SectionOptions []utils.SelectOption
	RoleOptions []utils.SelectOption
	EditorId string
	Name     string
	IsGroup  bool
	Text string
	Prefix string
	Suffix string
}

func (e *Engine) PromptBlockEditor(w io.Writer, ctx *PromptBlockEditor) error {
	return e.Template.ExecuteTemplate(w, "components/prompt_block_editor.html", ctx)
}
