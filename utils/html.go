package utils

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