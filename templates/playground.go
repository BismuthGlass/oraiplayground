package templates

import (
	"crow/oraiplayground/models"
	"log"
	"io"
	"slices"
)

func newPlaygroundBlockList(story *models.Story) blockEditorList {
	preset := story.ActivePreset()
	items := make([]blockEditorListItem, 0, len(preset.FavBlocks))
	for _, b := range story.PromptBlocks {
		if slices.Contains(preset.FavBlocks, b.Name) {
			items = append(items, blockEditorListItem{
				Name: b.Name,
				Enabled: slices.Contains(preset.EnabledBlocks, b.Name),
			})
		}
	}
	return blockEditorList{
		StoryName: story.Name,
		Items: items,
	}
}

func PlaygroundBlockList(w io.Writer, story *models.Story) error {
	ctx := newPlaygroundBlockList(story)
	err := engine.Template.ExecuteTemplate(w, "components/playground_block_list.html", ctx)
	if err != nil {
		log.Println(err)
	}
	return err
}
