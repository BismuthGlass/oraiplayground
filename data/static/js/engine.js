function setup() {
    let elemInInstruction = document.getElementById("input-instruction")
    let elemInPrompt = document.getElementById("input-prompt")
    let elemInContext = document.getElementById("input-context")
    let elemOutResponse = document.getElementById("output-response")
    let elemOutPrompt = document.getElementById("output-prompt")
    let elemBtGenerate = document.getElementById("bt-generate")

    function requestEnd(reqId) {
        data = {
            instruction: elemInInstruction.value,
            context: elemInContext.value,
            prompt: elemInPrompt.value,
        };
        fetch(`/prompt?id=${reqId}`, {
            method: "GET",
        })
        .then(response => response.json())
        .then(data => {
            elemOutPrompt.value = data.prompt;
            elemOutResponse.value = data.response;
        })
        .catch(err => {
            console.log("Error:", err)
        })
        .finally(() => {
            elemOutResponse.readOnly = false
            elemOutPrompt.readOnly = false
            elemBtGenerate.disabled = false
        })
    }

    function requestStart() {
        elemOutResponse.readOnly = true
        elemOutPrompt.readOnly = true
        elemBtGenerate.disabled = true

        let formData = new FormData()
        formData.append("instruction", elemInInstruction.value)
        formData.append("context", elemInContext.value)
        formData.append("prompt", elemInPrompt.value)
        fetch("/prompt", {
            method: "POST",
            body: formData,
        })
        .then(response => response.json())
        .then(data => {
            requestEnd(data.id)
        })
        .catch(err => {
            console.log("Error:", err)
        })
    }

    elemBtGenerate.addEventListener("click", requestStart)

    console.log("Environment loaded! Have fun :)")
}

document.addEventListener("DOMContentLoaded", setup);
