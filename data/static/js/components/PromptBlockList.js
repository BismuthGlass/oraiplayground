export class PromptBlockList {
    constructor(rootElem, clickHandler) {
        this.root = rootElem
        this.dragged = null
        this.form = rootElem.getElementsByTagName("form")[0]
        this.clickHandler = clickHandler

        this.setupRefreshListeners()
        this.setupEventListeners(rootElem.getElementsByClassName("block-list")[0])
    }
    
    setupEventListeners(listElem) {
        for (let e of listElem.children) {
            if (e.dataset.name) {
                e.getElementsByClassName("operations")[0].children[0].addEventListener("click", () => {
                    this.clickHandler(e.dataset.name)
                })
            }
            if (e.draggable) {
                e.addEventListener("dragstart", this.dragHandler.bind(this))
            }
            e.addEventListener("dragover", function(e) { e.preventDefault() })
            e.addEventListener("drop", this.dropHandler.bind(this))
        }
    }

    setupRefreshListeners() {
        htmx.on(this.root, "htmx:afterSettle", (e) => {
            let target = e.target
            if (target.id == "block-list")
                this.setupEventListeners(target)
        })
    }

    dragHandler(e) {
        this.dragged = e.target
    }

    dropHandler(e) {
        e.preventDefault()

        let dragTarget = e.currentTarget
        if (dragTarget.dataset.isGroup === "true" && this.dragged.dataset.isGroup === "true") {
            return
        }
        if (this.dragged == dragTarget) {
            return
        }

        this.form["moveTargetBlock"].value = this.dragged.dataset.name || ""
        this.form["moveDestinationBlock"].value = dragTarget.dataset.name || ""
        this.form["moveDestinationGroup"].value = dragTarget.dataset.parent || ""
        htmx.trigger(this.form, "prompt-block-list-update")
    }
}