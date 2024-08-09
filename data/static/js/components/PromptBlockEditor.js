class PromptBlockEditorList {
	constructor(root, storyName, deleteCallback, editCallback) {
		this.root = root
		this.storyName = storyName
		this.cbDelete = deleteCallback
		this.cbEdit = editCallback

		for (let row of root.children) {
			this.setupRow(row)
		}

		console.log("Setting up event listeners for the body")
		document.body.addEventListener("htmx:afterSettle", (e) => {
			let elem = e.detail.target
			console.log(elem)
		})
		document.body.addEventListener("htmx:afterSwap", (e) => {
			let elem = e.detail.target
			console.log(elem)
		})
		console.log(root)
		root.addEventListener("htmx:afterSettle", (e) => {
			let elem = e.detail.target
			console.log(elem)
		})
		root.addEventListener("htmx:afterSwap", (e) => {
			let elem = e.detail.target
			console.log(elem)
		})
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
		opFavorite.addEventListener("click", () => {
			console.log(`Favorited ${blockName}`)
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
	}

	isEdited() {
		return this.blockName.length > 0 && this.edited
	}

	edit(blockName) {
		// TODO: Add discard confirmation if there is already a block being edited
		htmx.ajax("GET", `/story/${this.storyName}/promptBlockEditor/${blockName}`, {
			target: this.root,
			swap: "outerHTML",
		})
	}
}

export class PromptBlockEditorMaster {
	constructor() {
		let list = document.getElementsByClassName("pblock-editor-table")[0]
		let editor = document.getElementsByClassName("pblock-editor")[0]
		this.list = new PromptBlockEditorList(list, "default", this.deleteBlock.bind(this), this.editBlock.bind(this))
		this.editor = new PromptBlockEditorEditor(editor)

		document.body.addEventListener("htmx:load", (e) => {
			let elt = e.detail.elt
			if (elt.classList.contains("pblock-editor-table")) {
				this.list = new PromptBlockEditorList(
					elt, 
					"default",				
					this.deleteBlock.bind(this),
					this.editBlock.bind(this),
				)
			} else if (elt.classList.contains("pblock-editor")) {
				this.editor = new PromptBlockEditor(elt)
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
	}

	editBlock(blockName) {
		console.log(`Edit called on ${blockName}`)
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
