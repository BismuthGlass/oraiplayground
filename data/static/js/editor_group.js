class EditorGroup {
    constructor(root) {
        this.root = root

        this.root.addEventListener("htmx:afterSwap", swapListener.bind(this))
    }

    swapListener(e) {
        console.log(e.detail)
    }
}
