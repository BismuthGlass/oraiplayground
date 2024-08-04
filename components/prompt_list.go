package components

import (
	"crow/oraiplayground/models"
	"crow/oraiplayground/templates"
)

type PromptList struct {
	DataRoot *models.PromptBlock
	Collapsed map[string]interface{}
}

func NewPromptList(root *models.PromptBlock) *PromptList {
	return &PromptList{
		DataRoot: root,
		Collapsed: make(map[string]interface{}),
	}
}

func (c *PromptList) Collapse(name string) {
	c.Collapsed[name] = nil
}

func (c *PromptList) Expand(name string) {
	delete(c.Collapsed, name)
}

func (c *PromptList) recurseTmplBlock(parent string, blocks []*models.PromptBlock, list []templates.PromptBlockListItem) []templates.PromptBlockListItem {
	for _, e := range blocks {
		entry := templates.PromptBlockListItem{
			Name: e.Name,
			Parent: parent,
			Active: e.Active,
			IsGroup: e.IsGroup,
			IsEmptyGroup: len(e.Children) == 0,
		}
		list = append(list, entry)
		list = c.recurseTmplBlock(e.Name, e.Children, list)
	}
	return list
}

func (c *PromptList) IntoTmplModel() templates.PromptBlockList {
	return templates.PromptBlockList{
		Items: c.recurseTmplBlock("", c.DataRoot.Children, nil),
	}
}
