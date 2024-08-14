class BlockEditorList {
	constructor(root, storyName, deleteCallback, editCallback) {
		this.root = root
		this.storyName = storyName
		this.cbDelete = deleteCallback
		this.cbEdit = editCallback

		for (let row of root.children) {
			this.setupRow(row)
		}

		window.blockEditor.setupRow = this.setupRow.bind(this)
	}

	setupRow(row) {
		let blockName = row.dataset.blockName
		let opFavorite = row.getElementsByClassName("op-favorite")[0]
		let opDelete = row.getElementsByClassName("op-delete")[0]
		let opEdit = row.getElementsByClassName("op-edit")[0]

		opEdit.addEventListener("click", () => {
			this.cbEdit(blockName)
		})
	}
}

class BlockEditorForm {
	constructor(root, storyName) {
		this.root = root
		this.storyName = storyName
		this.blockName = root.dataset.blockName
		this.edited = false

		this.saveBt = root.getElementsByClassName("save-bt")[0]
		this.discardBt = root.getElementsByClassName("discard-bt")[0]
		this.closeBt = root.getElementsByClassName("close-bt")[0]

		let roleSelect = root.elements["role"]
		let text = root.elements["text"]
		let name = root.elements["name"]
		let compiled = root.elements["compiled"]

		roleSelect.addEventListener("change", this.whenEdited.bind(this))
		text.addEventListener("keydown", this.whenEdited.bind(this))
		name.addEventListener("keydown", this.whenEdited.bind(this))
		compiled.addEventListener("change", this.whenEdited.bind(this))
	}

	whenEdited() {
		this.edited = true
		this.saveBt.disabled = false
		this.discardBt.disabled = false
		this.closeBt.disabled = true
	}

	isEdited() {
		return this.edited
	}

	edit(blockName) {
		htmx.ajax("GET", `/story/${this.storyName}/blockEditor/edit/${blockName}`, {
			target: this.root,
			swap: "outerHTML",
		})
	}
}

export class BlockEditorMaster {
	constructor() {
		window.blockEditor = {}
		window.blockEditor.setupList = (e) => {
			this.list = new BlockEditorList(e, "default", this.deleteBlock.bind(this), this.editBlock.bind(this))
		}
		window.blockEditor.confirmDelete = (elt, blockName) => {
			if (this.editor.blockName == blockName) {
				alert("Block is open for editing")
			}
			if (confirm(`Delete ${blockName}?`)) {
				htmx.trigger(elt, 'confirmed')
			}
		}

		let list = document.getElementsByClassName("pblock-editor-table")[0]
		let editor = document.getElementsByClassName("pblock-editor")[0]
		this.list = new BlockEditorList(list, "default", this.deleteBlock.bind(this), this.editBlock.bind(this))
		this.editor = new BlockEditorForm(editor, "default")

		document.body.addEventListener("htmx:load", (e) => {
			let elt = e.detail.elt
			if (elt.classList.contains("pblock-editor")) {
				this.editor = new BlockEditorForm(elt, "default")
			}
		})
	}

	deleteBlock(blockName) {
		if (this.editor.blockName == blockName) {
			alert("Block is open for editing")
			return
		}
		if (!confirm(`Delete ${blockName}?`)) {
			return
		}
		// TODO: Delete block here
	}

	editBlock(blockName) {
		if (this.editor.isEdited()) {
			if (confirm("Discard changes?")) {
				this.editor.edit(blockName)
			}
		} else {
			this.editor.edit(blockName)
		}
	}
}

