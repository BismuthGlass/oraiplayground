import { PromptBlockList } from "./components/PromptBlockList.js"
import { TabController, setupTabs } from "./components/TabController.js"
import { EditorLayout } from "./components/EditorLayout.js"
import { setupBlockEditor } from "./components/BlockEditor.js"
import { setupPlayground } from "./components/Playground.js"
import { setupAiResponseControls } from "./components/OutputController.js"

function setup() {
	window.setupTabs = setupTabs

	new TabController(
		document.getElementById("navigation").children,
		document.getElementById("full-window-div").children,
		"story",
	)

	setupBlockEditor(document.getElementsByClassName("tab-prompt")[0])
	setupPlayground(document.getElementsByClassName("tab-story")[0])
	setupAiResponseControls(document.getElementById("ai-response"))

	console.log("setup done!")

	if (true) {
		htmx.logger = function(elt, event, data) {
			if(console) {
				console.log("INFO:", event, elt, data)
			}
		}
	}
}

document.addEventListener("DOMContentLoaded", setup)

