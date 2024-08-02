import { TabController } from "./TabController.js"

export class OutputController {
    constructor(root) {
        this.root = root
        this.resultElem = document.getElementById("ai-response-response")
        this.promptElem = document.getElementById("ai-response-prompt")
        this.promptPreviewElem = document.getElementById("ai-response-prompt-preview")

        this.requestId = null

        this.btCancel = document.getElementById("ai-response-bt-cancel")
        this.btGenerate = document.getElementById("ai-response-bt-generate")
        this.btContinue = document.getElementById("ai-response-bt-continue")
        this.btPreviewPrompt = document.getElementById("ai-response-bt-preview")

        new TabController(
            root.getElementsByClassName("button-group")[0].children,
            root.children,
            "response"
        )

        this.btGenerate.addEventListener("click", this.requestGenerate.bind(this))
        this.btContinue.addEventListener("click", this.requestContinue.bind(this))
        this.btCancel.addEventListener("click", this.cancelPendingRequest.bind(this))
        this.btCancel.disabled = true
        this.btPreviewPrompt.addEventListener("click", this.promptPreviewRequest.bind(this))
    }

    issueRequest(dataHandler, continuationPrompt) {
        this.btGenerate.disabled = true
        this.btContinue.disabled = true

        let data = new FormData()
        if (continuationPrompt) {
            data.append("continue", continuationPrompt)
        }
        
        //let data = new FormData()
        fetch("/promptRequest", {
            method: "POST",
            body: data, 
        })
        .then(value => value.json())
        .then((data) => {
            this.requestId = data.id
            this.btCancel.disabled = false
            fetch(`/promptRequest?id=${data.id}`, {
                method: "GET",
            })
            .then(data => data.json())
            .then((data) => {
                dataHandler(data)
            })
        })
    }

    requestGenerate() {
        this.issueRequest((data) => {
            this.resultElem.value = data.response
            this.promptElem.value = data.prompt
            this.resetButtons()
        })
    }

    requestContinue() {
        this.issueRequest((data) => {
            this.resultElem.value += data.response
            this.promptElem.value = data.prompt
            this.resetButtons()
        }, this.resultElem.value)
    }

    cancelPendingRequest() {
        if (this.requestId) {
            fetch(`/promptRequestDelete?id=${this.requestId}`, {
                method:"DELETE"
            })
        }
    }

    resetButtons() {
        this.btCancel.disabled = true
        this.btContinue.disabled = false
        this.btGenerate.disabled = false
    }

    promptPreviewRequest() {
        fetch("/promptPreview", {
            method: "GET"
        })
        .then((value) => value.json())
        .then((data) => {
            this.promptPreviewElem.value = data["prompt"]
        })
    }
}