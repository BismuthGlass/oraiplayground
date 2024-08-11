class BlockEditorList {
	constructor(root, storyName, deleteCallback, editCallback) {
		this.root = root
		this.storyName = storyName
		this.cbDelete = deleteCallback
		this.cbEdit = editCallback

		for (let row of root.children) {
			this.setupRow(row)
		}
	}

	setupRow(row) {
		let blockName = row.dataset.blockName
		let opFavorite = row.getElementsByClassName("op-favorite")[0]
		let opDelete = row.getElementsByClassName("op-delete")[0]
		let opEdit = row.getElementsByClassName("op-edit")[0]

		opEdit.addEventListener("click", () => {
			this.cbEdit(blockName)
		})
		opDelete.addEventListener("click", () => {
			this.cbDelete(blockName)
		})
	}
}

class BlockEditorForm {
	constructor(root, storyName) {
		this.root = root
		this.storyName = storyName
		this.blockName = root.dataset.blockName
		this.edited = false
		console.log(this.storyName)

		this.saveBt = root.getElementsByClassName("save-bt")[0]
		this.discardBt = root.getElementsByClassName("discard-bt")[0]
		this.closeBt = root.getElementsByClassName("close-bt")[0]

		let roleSelect = root.elements["role"]
		let text = root.elements["text"]
		let name = root.elements["name"]

		roleSelect.addEventListener("change", this.whenEdited.bind(this))
		text.addEventListener("keypress", this.whenEdited.bind(this))
		name.addEventListener("keypress", this.whenEdited.bind(this))
	}

	whenEdited() {
		console.log("Edited!")
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
		let list = document.getElementsByClassName("pblock-editor-table")[0]
		let editor = document.getElementsByClassName("pblock-editor")[0]
		this.list = new BlockEditorList(list, "default", this.deleteBlock.bind(this), this.editBlock.bind(this))
		this.editor = new BlockEditorForm(editor, "default")

		document.body.addEventListener("htmx:load", (e) => {
			let elt = e.detail.elt
			if (elt.classList.contains("pblock-editor-table")) {
				this.list = new BlockEditorList(
					elt, 
					"default",				
					this.deleteBlock.bind(this),
					this.editBlock.bind(this),
				)
			} else if (elt.classList.contains("pblock-editor")) {
				console.log("Setting up form")
				this.editor = new BlockEditorForm(elt, "default")
			}
		})

		document.body.addEventListener("evtRefreshBlockList", (e) => {
			if (e.details.targets.includes("pblock-editor-table")) {
				htmx.ajax("GET", "/story/default/promptBlockEditor/list", {
					target: this.list.root,
					swap: "outerHTML",
				})
			}
		});
	}

	deleteBlock(blockName) {
		console.log(`Delete called on ${blockName}`)
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
		console.log(`Edit called on ${blockName}`)
		if (this.editor.isEdited()) {
			if (confirm("Discard changes?")) {
				this.editor.edit(blockName)
			}
		} else {
			this.editor.edit(blockName)
		}
	}
}

export class PromptBlockEditor {
    constructor(root) {
        this.root = root
        this.empty = !root.dataset.blockName || root.dataset.blockName === ""

        this.root.dataset.edited = "false"

        if (!this.empty) {
            this.saveBt = root.getElementsByClassName("save-bt")[0]
            this.discardBt = root.getElementsByClassName("discard-bt")[0]
            this.closeBt = root.getElementsByClassName("close-bt")[0]
            this.modifiedIndicator = root.getElementsByClassName("modified-indicator")[0]

            root.elements["section"].addEventListener("change", this.changeHandler.bind(this))
            if (root.elements["prefix"])
                root.elements["prefix"].addEventListener("input", this.changeHandler.bind(this))
            if (root.elements["suffix"])
                root.elements["suffix"].addEventListener("input", this.changeHandler.bind(this))
            if (root.elements["text"])
                root.elements["text"].addEventListener("input", this.changeHandler.bind(this))

            this.saveBt.addEventListener("click", this.saveHandler.bind(this))
            this.discardBt.addEventListener("htmx:confirm", this.confirmHandle.bind(this))
        }
    }

    changeHandler(e) {
        this.saveBt.disabled = false
        this.discardBt.disabled = false
        this.closeBt.disabled = true
        this.root.dataset.edited = "true"
        this.modifiedIndicator.innerHTML = "*"
    }

    saveHandler(e) {
        this.saveBt.disabled = true
        this.discardBt.disabled = true
        this.closeBt.disabled = false
        this.modifiedIndicator.innerHTML = ""
        this.root.dataset.edited = "false"
        htmx.trigger(this.root, `${this.root.id}-save`)
    }

    confirmHandle(e) {
        e.preventDefault()
        if (!this.isModified) {
            e.detail.issueRequest(true)
        } else {
            let ok = confirm(e.detail.question)
            if (ok) {
                e.detail.issueRequest(true)
            }
        }
    }

    isModified() {
        return this.modifiedIndicator.innerHTML == "*"
    }
}
