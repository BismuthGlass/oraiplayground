import { PromptBlockList } from "./components/PromptBlockList.js"
import { TabController } from "./components/TabController.js"
import { EditorLayout } from "./components/EditorLayout.js"
import { BlockEditorMaster } from "./components/BlockEditor.js"
import { setupPromptInfoPopup } from "./components/PromptInfo.js"

class StoryMasterController {
	constructor(root) {
		new TabController(
			document.getElementById("navigation").children,
			document.getElementById("full-window-div").children,
			"story",
		)

		this.editorLayout = new EditorLayout(
			document.getElementsByClassName("story-layout")[0],
		)

		new PromptBlockList(
			document.getElementById("prompt-block-list"),
			this.openBlock.bind(this),
		)
	}

	openBlock(name) {
		this.editorLayout.open(name)
	}
}

function setup() {
	new TabController(
		document.getElementById("navigation").children,
		document.getElementById("full-window-div").children,
		"story",
	)

	new PromptBlockList(
		document.getElementById("prompt-block-list"),
		function(name) {
			editorLayout.open(name)
		},
	)

	new BlockEditorMaster()

	setupPromptInfoPopup()

	console.log("setup done!")
}

document.addEventListener("DOMContentLoaded", setup)

