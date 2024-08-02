import { PromptBlockEditor } from "./PromptBlockEditor.js";
import { OutputController } from "./OutputController.js";

export class EditorLayout {
    constructor(root) {
        this.root = root;
        this.editors = [];

        new OutputController(document.getElementById("ai-response"))
        
        htmx.on(this.root, "htmx:afterSwap", (e) => {
            let target = e.target
            if (target.classList.contains("prompt-block-editor")) {
                new PromptBlockEditor(target)
            }
        })
    }

    openEditor(editorId, block) {
        htmx.ajax("GET", `/blockEditor?eid=${editorId}&block=${block}`, {target: `#${editorId}`, swap: "outerHTML"})
    }

    open(name) {
        let editors = this.root.getElementsByClassName("prompt-block-editor")
        // Check if already open
        for (let editor of editors) {
            if (editor.dataset.blockName == name) {
                return
            }
        }
        // Check for empty editor
        for (let editor of editors) {
            if (editor.dataset.blockName == "") {
                this.openEditor(editor.id, name)
                return
            }
        }
        // Use first editor if not edited
        if (editors[0].dataset.edited == "false") {
            this.openEditor(editors[0].id, name)
        }
    }
}
