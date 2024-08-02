export class TabController {
    constructor(buttons, tabs, defaultTab) {
        this.tabs = {}
        
        for (let tab of tabs) {
            if (!tab.dataset.tabName) {
                continue
            }
            this.tabs[tab.dataset.tabName] = tab
        }

        for (let btn of buttons) {
            if (!btn.dataset.tabName) {
                continue
            }
            btn.addEventListener("click", () => {
                this.activate(btn.dataset.tabName)
            })
        }

        this.activate(defaultTab)
    }

    activate = (tab_name) => {
        for (let tab of Object.entries(this.tabs)) {
            tab[1].style.display = "none"
        }
        this.tabs[tab_name].style.display = ""
    }
}