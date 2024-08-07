class PromptBlockEditorList {
	constructor(root, storyName, deleteCallback, editCallback) {
		this.root = root
		this.storyName = storyName
		this.cbDelete = deleteCallback
		this.cbEdit = editCallback
	}

	setRoot(root) {
		this.root.addEventListener("htmx:afterSettle", function() {
			console.log("Detected settle on block list")
			// setRoot here...
		})
	}

	setupRow(row) {
		let blockName = row.data.blockName
		let opFavorite = row.getElementsByClassName("op-favorite")[0]
		let opDelete = row.getElementsByClassName("op-delete")[0]
		let opEdit = row.getElementsByClassName("op-edit")[0]

		opEdit.addEventListener("click", () => {
			this.cbEdit(blockName)
		})
		opDelete.addEventListener("click", () => {
			this.cbDelete(blockName)
		})
		opFavorite.addEventListener("click", () => {
			htmx.ajax("PUT", `/story/${this.storyName}/promptBlock/${blockName}/favorite`, {
				target: this.root,
				swap: "outerHTML",
			})
		})
	}
}

class PromptBlockEditorEditor {
	constructor(root, storyName) {
		this.storyName = storyName
		this.blockName = ""
		this.edited = false

		this.setRoot(root)
	}

	setRoot(root) {
		this.root.addEventListener("htmx:afterSettle", function() {
			console.log("Detected settle on block editor")
			// setRoot here...
		})
	}

	isEdited() {
		return this.blockName.length > 0 && this.edited
	}

	edit(blockName) {
		// TODO: Add discard confirmation if there is already a block being edited
		htmx.ajax("GET", `/story/${this.storyName}/promptBlock/${blockName}/editor`, {
			target: this.root,
			swap: "outerHTML",
		})
	}
}

export class PromptBlockEditorMaster {
	constructor() {
		let list = document.getElementsByClassName("pblock-editor-table")[0]
		let editor = document.getElementsByClassName("pblock-editor")[0]
		new PromptBlockEditorList(list, this.deleteBlock.bind(this), this.editBlock.bind(this))
		new PromptBlockEditorEditor(editor)
	}

	deleteBlock(blockName) {

	}

	editBlock(blockName) {
		
	}
}
