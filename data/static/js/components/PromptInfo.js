import { TabController } from "./TabController.js"

export function setupPromptInfoPopup() {
	let storyName = document.body.dataset.storyName
	let button = document.getElementById("open-prompt-info")
	let popup = document.getElementById("prompt-info")
	if (!button || !popup) {
		return
	}

	function setupTabs(contentRoot) {
		let nav = contentRoot.querySelector("nav")
		new TabController(nav.children, contentRoot.children, "next-prompt")
	}

	button.addEventListener("click", function() {
		popup.classList.remove("hidden")
		htmx.ajax("GET", `/story/${storyName}/promptInfo`, {
			target: "#prompt-info .content",
			swap: "outerHTML"
		})
	})

	document.body.addEventListener("htmx:afterSwap", function(e) {
		if (e.target.classList.contains("prompt-info-content")) {
			setupTabs(e.target)
		}
	})

	let background = popup.getElementsByClassName("background")[0]
	background.addEventListener("click", function() {
		popup.classList.add("hidden")
	})
}

