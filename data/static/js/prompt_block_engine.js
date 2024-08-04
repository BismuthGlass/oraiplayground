

export class Block {
    constructor(isGroup, name) {
        this.isGroup = isGroup
        this.name = name
        this.enabled = false
        this.children = []
    }

    addChild(child) {
        this.children.push(child)
    }
}

export function seedData() {
    let groupBlocks = {
        "group_1": new Block(true, "group_1"),
        "group_2": new Block(true, "group_2"),
        "group_3": new Block(true, "group_3"),
    }

    let childBlocks = {
        "block_1": new Block(false, "block_1"),
        "block_2": new Block(false, "block_2"),
        "block_3": new Block(false, "block_3"),
        "block_4": new Block(false, "block_4"),
        "block_5": new Block(false, "block_5"),
        "block_6": new Block(false, "block_6"),
    }

    groupBlocks["group_1"].addChild(childBlocks["block_1"])
    groupBlocks["group_2"].addChild(childBlocks["block_2"])
    groupBlocks["group_2"].addChild(childBlocks["block_3"])

    return [
        childBlocks["block_5"],
        groupBlocks["group_1"],
        groupBlocks["group_3"],
        childBlocks["block_4"],
        childBlocks["block_6"],
        groupBlocks["group_2"],
        childBlocks["block_5"],
        groupBlocks["group_1"],
        groupBlocks["group_3"],
        childBlocks["block_4"],
        childBlocks["block_6"],
        groupBlocks["group_2"],
    ]
}

function findBlock(data, name) {
    for (let b of data) {
        if (b.name == name) {
            return b
        }
        if (b.isGroup) {
            let found = findBlock(b.children, name)
            if (found != null) {
                return found
            }
        }
    }
    return null
}

function getBlockIndex(data, name) {
    for (let i = 0; i < data.length; i++) {
        if (data[i].name == name) {
            return i
        }
    }
    return -1
}

function findGroup(data, name) {
    for (let b of data) {
        if (!b.isGroup) {
            continue
        }
        if (b.name == name) {
            return b
        }
        if (b.isGroup) {
            let found = findBlock(b.children, name)
            if (found != null) {
                return found
            }
        }
    }
    return null
}

export function renderPromptBlockList(parentElem, data) {
    function renderBlock(data) {
        let elem = document.createElement("div")
        let elemName = document.createElement("p")
        elem.appendChild(elemName)
        elemName.innerHTML = data.name
        elem.dataset.isGroup = false
        elem.dataset.name = data.name
        elem.draggable = true
        return elem
    }

    function renderGroup(data) {
        let elem = document.createElement("div")
        let elemName = document.createElement("p")
        elem.appendChild(elemName)
        elemName.innerHTML = data.name
        elem.dataset.isGroup = true
        elem.dataset.name = data.name
        elem.draggable = true
        return elem
    }

    function renderGroupChild(data, parent) {
        let elem = document.createElement("div")
        let elemName = document.createElement("p")
        elem.appendChild(elemName)
        elemName.innerHTML = data.name
        elem.classList.add("in-group")
        elem.dataset.isGroup = false
        elem.dataset.name = data.name
        elem.dataset.parent = parent.name
        elem.draggable = true
        return elem
    }

    function renderGroupPlaceholder(parent) {
        let elem = document.createElement("div")
        let elemName = document.createElement("p")
        elem.appendChild(elemName)
        elemName.innerHTML = "Empty"
        elem.classList.add("in-group")
        elem.classList.add("empty-group")
        elem.dataset.parent = parent.name
        return elem
    }

    parentElem.innerHTML = ""

    console.log(data)
    for (let block of data) {
        if (block.isGroup) {
            parentElem.appendChild(renderGroup(block))
            if (block.children.length > 0) {
                for (let child of block.children) {
                    parentElem.appendChild(renderGroupChild(child, block))
                }
            } else {
                parentElem.appendChild(renderGroupPlaceholder(block))
            }
        } else {
            parentElem.appendChild(renderBlock(block))
        }
    }

    // Set up dragging
    let currentDrag = null

    for (let block of parentElem.children) {
        block.addEventListener("dragstart", function(e) {
            currentDrag = this
            console.log("Dragging...")
        })

        block.addEventListener('dragover', (e) => {
            e.preventDefault();
        });

        block.addEventListener("drop", function(e) {
            if (this.dataset.name == currentDrag.dataset.name) {
                return
            }

            // Remove element from parent
            let draggedGroup = findGroup(data, currentDrag.dataset.parent)
            let draggedHomeArray 
            if (draggedGroup) {
                draggedHomeArray = draggedGroup.children
            } else {
                draggedHomeArray = data
            }
            let index = getBlockIndex(draggedHomeArray, currentDrag.dataset.name)
            let draggedBlock = draggedHomeArray[index]
            draggedHomeArray.splice(index, 1)

            // Find destination
            let group = findGroup(data, this.dataset.parent)
            let destination;
            if (group) {
                destination = group.children
            } else {
                destination = data
            }
            if (!block) {
                destination.add(draggedBlock)
            } else {
                let blockIndex = getBlockIndex(destination, this.dataset.name)
                destination.splice(blockIndex, 0, draggedBlock)
            }
            
            renderPromptBlockList(parentElem, data)
        })
    }
}

document.addEventListener("DOMContentLoaded", function() {
    let elem = document.getElementById("block-list")
    renderPromptBlockList(elem, seedData())
})
