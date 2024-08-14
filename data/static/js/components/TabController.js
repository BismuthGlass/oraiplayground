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

export function setupTabs(root, defaultTab) {
	let nav = null
	for (let c of root.children) {
		if (c.nodeName == "NAV") {
			nav = c
			break
		}
	}
	if (nav == null) {
		return
	}

	let tabs = {}
	let activate = function(tab_name) {
		for (let tab of Object.entries(tabs)) {
			tab[1].style.display = "none"
		}
		tabs[tab_name].style.display = ""
	}

	for (let tab of root.children) {
		if (!tab.dataset.tabName) {
			continue
		}
		tabs[tab.dataset.tabName] = tab
	}

	for (let btn of nav.children) {
		if (!btn.dataset.tabName) {
			continue
		}
		btn.addEventListener("click", function() {
			activate(btn.dataset.tabName)
		})
	}

	activate(defaultTab)
}

