import { TabController } from "./TabController.js"

const RequestState = {
	Ready: 0,
	Pending: 1,
	Canceled: 2,
}

function requestAiGen(storyName, cue) {
	return fetch(`/story/${storyName}/gen`, {
		method: "POST",
		body: JSON.stringify({
			cue: cue
		})
	})
		.then((r) => r.json())
}

export function setupAiResponseControls(root) {
	let storyName = document.body.dataset.storyName

	let requestState = RequestState.Ready
	let requestId = 0
	let targetTextarea = null

	let generateBt = root.getElementsByClassName("generate-bt")[0]
	let continueBt = root.getElementsByClassName("continue-bt")[0]
	let cancelBt = root.getElementsByClassName("cancel-bt")[0]
	let responseArea = root.getElementsByClassName("response")[0]

	function issueRequest(target, append) {
		requestState = RequestState.Pending
		generateBt.disabled = true
		continueBt.disabled = true
		targetTextarea = target

		requestAiGen(storyName, "")
			.then(async (body) => {
				requestId = body.id
				cancelBt.disabled = false
				await fetch(`/story/${storyName}/gen/${body.id}`, {
					method: "GET",
				})
					.then((r) => r.json())
					.then((body) => {
						console.log(body)
						if (body.err) {
							return
						}
						if (append) {
							target.innerHTML += body.response
						} else {
							target.innerHTML = body.response
						}
					})
			})
			.finally(() => {
				generateBt.disabled = false
				continueBt.disabled = false
				cancelBt.disabled = true
			})
	}

	function issueCancel() {
		requestState = RequestState.Canceled
		cancelBt.enabled = false
		fetch(`/story/${storyName}/gen/${requestId}`, {
			method: "DELETE",
		})
	}

	generateBt.addEventListener("click", function() { issueRequest(responseArea, false) })
	continueBt.addEventListener("click", function() { issueRequest(responseArea, true) })
	cancelBt.addEventListener("click", function() { issueCancel() })
}

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
        fetch("/aiRequest", {
            method: "POST",
            body: data, 
        })
        .then(value => value.json())
        .then((data) => {
            this.requestId = data.id
            this.btCancel.disabled = false
            fetch(`/aiRequest?id=${data.id}`, {
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
            fetch(`/aiRequest?id=${this.requestId}`, {
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
