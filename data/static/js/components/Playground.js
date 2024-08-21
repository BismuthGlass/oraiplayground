export function setupPlayground(tabRoot) {
	let storyName = document.body.dataset.storyName

	function openEditor(blockName) {
		function open(elt, blockName) {
			htmx.ajax("GET", `/story/${storyName}/blockEditor/edit/${blockName}`, {
				target: elt,
				swap: "outerHTML",
			})
		}

		let editors = tabRoot.getElementsByClassName("pblock-editor")
		for (let e of editors) {
			if (e.dataset.blockName == "") {
				open(e, blockName)
				return
			} else if (e.dataset.blockName == blockName) {
				return
			}
		}
		for (let e of editors) {
			if (e.dataset.edited != "true") {
				open(e, blockName)
				return
			}
		}
	}

	function closeEditor(blockName) {
		let editors = tabRoot.getElementsByClassName("pblock-editor")
		for (let e of editors) {
			if (e.dataset.blockName == blockName) {
				let form = document.createElement('form');
				form.className = 'v-flex pblock-editor';
				form.setAttribute('data-block-name', '');
				e.replaceWith(form)
				return
			}
		}
	}

	window.playground = {}
	window.playground.openEditor = openEditor
	window.playground.closeEditor = closeEditor
}

