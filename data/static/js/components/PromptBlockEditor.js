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