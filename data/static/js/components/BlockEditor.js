function setupList(list) {
	let storyName = document.body.dataset.storyName
	for (let i of list.children) {
		i.addEventListener("dragstart", function(e) {
			e.dataTransfer.clearData()
			e.dataTransfer.setData("text/plain", i.dataset.blockName)
		})
		i.addEventListener("dragover", function(e) { e.preventDefault() })
		i.addEventListener("drop", function(e) {
			e.preventDefault()

			let dragTarget = e.currentTarget.dataset.blockName
			let dragged = e.dataTransfer.getData("text")
			if (dragged == dragTarget) {
				return
			}
			htmx.ajax("put", `/story/${storyName}/blockEditor/move/${dragged}/${dragTarget}`, { swap: "none" })
		})
	}
}

function setupEditor(form) {
	let storyName = document.body.dataset.storyName

	let saveBt = form.getElementsByClassName("save-bt")[0]
	let discardBt = form.getElementsByClassName("discard-bt")[0]
	let closeBt = form.getElementsByClassName("close-bt")[0]
	function whenEdited() {
		form.dataset.edited = "true"
		saveBt.disabled = false
		discardBt.disabled = false
		closeBt.disabled = true
	}
	let roleSelect = form.elements["role"]
	let text = form.elements["text"]
	let name = form.elements["name"]
	let compiled = form.elements["compiled"]
	roleSelect.addEventListener("change", whenEdited)
	text.addEventListener("keydown", whenEdited)
	name.addEventListener("keydown", whenEdited)
	compiled.addEventListener("change", whenEdited)
}

export function setupBlockEditor(tabRoot) {
	let storyName = document.body.dataset.storyName

	function openEditor(elt, blockName) {
		htmx.ajax("GET", `/story/${storyName}/blockEditor/edit/${blockName}`, {
			target: elt,
			swap: "outerHTML",
		})
	}

	window.blockEditor = {}
	window.blockEditor.setupList = setupList
	window.blockEditor.confirmDelete = function(elt, blockName) {
		let editors = document.getElementsById("pblock-editor")
		for (let e of editors) {
			if (e.dataset.blockName == blockName) {
				alert("Block is open for editing")
				return
			}
		}
		if (confirm(`Delete ${blockName}?`)) {
			htmx.trigger(elt, 'confirmed')
		}
	}
	window.blockEditor.openEditor = function(blockName) {
		let editors = tabRoot.getElementsByClassName("pblock-editor")
		for (let e of editors) {
			if (e.dataset.blockName == "") {
				openEditor(e, blockName)
				return
			} else if (e.dataset.blockName == blockName) {
				return
			}
		}
		for (let e of editors) {
			if (e.dataset.edited != "true") {
				openEditor(e, blockName)
				return
			}
		}
	}
	window.blockEditor.setupEditor = setupEditor

	setupList(tabRoot.getElementsByClassName("pblock-editor-table")[0])
	let editors = tabRoot.getElementsByClassName("pblock-editor")
	for (let e of editors) {
		setupEditor(e)
	}
}

