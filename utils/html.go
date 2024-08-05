package utils

import (
	"bytes"
)

type Select struct {
	Options []SelectOption
	Value string
}

type SelectOption struct {
	Value    string
	Name     string
	Selected bool
}

func SetSelection(options []SelectOption, value string) {
	for i := 0; i < len(options); i++ {
		if options[i].Value == value {
			options[i].Selected = true
		} else {
			options[i].Selected = false
		}
	}
}

func TFunRenderSelectOptions(s *Select) string {
	var buffer bytes.Buffer
	return buffer.String()
}
