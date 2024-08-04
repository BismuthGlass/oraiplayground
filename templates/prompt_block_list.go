package templates

import (
	"io"
	"log"
)

type PromptBlockListItem struct {
	Name         string
	Parent       string
	IsGroup      bool
	IsEmptyGroup bool
	IsCollapsed  bool
	Active       bool
}

type PromptBlockList struct {
	Items []PromptBlockListItem
}

func (e *Engine) PromptBlockList(w io.Writer, ctx *PromptBlockList) error {
	err := e.Template.ExecuteTemplate(w, "components/prompt_block_list.html", ctx)
	if err != nil {
		log.Println(err)
	}
	return err
}