package models

import (
	"slices"
	"log"
)

type PromptSettings struct {
	EnabledBlocks []string
	VarOverrides []StoryVariable
	FavBlocks []string
	FavVars []string
}

func (s *PromptSettings) enableBlock(name string) {
	if !slices.Contains(s.EnabledBlocks, name) {
		s.EnabledBlocks = append(s.EnabledBlocks, name);
	}
}

func (s *PromptSettings) disableBlock(name string) {
	s.EnabledBlocks = slices.DeleteFunc(s.EnabledBlocks, func(n string) bool { return n == name }) 
}

func (s *PromptSettings) toggleBlockFavorite(name string) {
	if slices.Contains(s.FavBlocks, name) {
		s.FavBlocks = slices.DeleteFunc(s.FavBlocks, func(n string) bool { return n == name })
	} else {
		s.FavBlocks = append(s.FavBlocks, name)
	}
}
