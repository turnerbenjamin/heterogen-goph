
let tableConfig


//* QUERY PARAMS
function getColumnParams() {
    if (!tableConfig) return ""
    return tableConfig.columns.join(",")
}

function getSortParams() {
    return tableConfig?.sort ?? ""
}

//* SORT HANDLER
function handleSorting(e){
    const buttonEl = e.target.closest("button")
    const sortDirection = buttonEl.dataset.sort
    const fieldLabel = buttonEl.closest("th").dataset.label

    if(sortDirection == "") tableConfig.sort = `${fieldLabel}`
    if(sortDirection == "ASC") tableConfig.sort = `-${fieldLabel}`
    if(sortDirection == "DESC") tableConfig.sort = ""
    htmx.trigger("#table-container", "refresh-table")
}




//* INITIAL DATA FETCH
function loadTableConfig() {
    storedConfig = localStorage.getItem("userTableConfig")
    if (!storedConfig) return;
    try {
        tableConfig = JSON.parse(storedConfig);
    } catch (error) {
        console.error("Error parsing stored config:", error);
    }
}

document.addEventListener('DOMContentLoaded', () => {
    loadTableConfig()
    htmx.trigger("#table-container", "refresh-table")
});


//* ON TABLE REFRESH

//Initialises the table config where no stored config found
function initialiseTableConfig() {
    tableConfig = {
        columns: [],
        columnConfig: {}
    }

    const headers = document.querySelectorAll("th")
    headers.forEach(header => {
        headerName = header.innerText;
        headerConfig = {
            width: header.getBoundingClientRect().width
        }
        tableConfig.columns.push(headerName);
        tableConfig.columnConfig[headerName] = headerConfig
    })
    localStorage.setItem("userTableConfig", JSON.stringify(tableConfig))
}

//Apply config after refresh to keep column widths correct
function applyConfig() {
    const headers = document.querySelectorAll("th")
    for (const header of headers) {
        const headerName = header.innerText;
        const headerConfig = tableConfig.columnConfig[headerName]
        header.style.width = headerConfig.width + "px"
    }
    headers.forEach(header => {
        headerName = header.innerText;
        headerConfig = {
            width: header.getBoundingClientRect().width
        }
    })
}

//ATTACH HANDLERS TO MANAGE COL RESIZE AND REORDER
function addTableHandlers() {
    const table = document.querySelector("table");
    if (!table) return

    //Table resizers
    const tableResizers = table.querySelectorAll(".resizer")
    const tableHeight = table.getBoundingClientRect().height;
    for (const tableResizer of tableResizers) {
        tableResizer.style.height = tableHeight + "px"
        applyDragHandler(tableResizer, handleResize, "resizing")
    }
    handleReorder()
}

function handleTableRefresh() {
    if (!tableConfig) initialiseTableConfig()
    else applyConfig()
    addTableHandlers()
}

//Update config where column sizes changed
function updateColumnSizes() {
    const headers = document.querySelectorAll("th")
    headers.forEach(header => {
        const headerName = header.innerText;
        const headerConfig = tableConfig.columnConfig[headerName]
        const width = header.getBoundingClientRect().width;
        headerConfig.width = width
    })
    localStorage.setItem("userTableConfig", JSON.stringify(tableConfig))
}


function handleResize(target, e) {
    const resizer = target
    const header = resizer.closest("th")
    const { width, right } = header.getBoundingClientRect();
    const dx = e.clientX - right;
    const newWidth = width + dx
    header.style.width = newWidth + 'px';
}


let moving
let movingTo
function handleReorder() {

    const handleDragStart = (e) => {
        moving = e.target.closest("th")
        movingTo = null
    }
    const handleDragEnd = (e) => {
        if(!movingTo) return
        handleSwapColumns(moving, movingTo)
    }

    const handleDragOver = (e) => {
        e.preventDefault()

        if (movingTo) movingTo.classList.remove("dragged-over");

        if (!moving) return
        const target = e.target.closest("th")
        if (!target || target == moving) return;

        movingTo = target;
        movingTo.classList.add("dragged-over")
    }

    const table = document.querySelector("table")
    table.addEventListener("dragstart", handleDragStart)
    table.addEventListener("dragend", handleDragEnd)
    table.addEventListener("dragover", handleDragOver)

}

const handleSwapColumns = (move, moveTo) => {
    const updatedColumns = [];
    let colFound = false
    for (let column of tableConfig.columns) {
        if (column == moveTo.innerText && !colFound) updatedColumns.push(move.innerText)
        if (column != move.innerText) {
            updatedColumns.push(column)
        }
        else {
            colFound = true
        }

        if (column == moveTo.innerText && colFound) updatedColumns.push(move.innerText)
    }
    tableConfig.columns = updatedColumns
    console.log(tableConfig)
    localStorage.setItem("userTableConfig", JSON.stringify(tableConfig))
    htmx.trigger("#table-container", "refresh-table")
}


//utility function - used to handle and clean-up drag listeners
function applyDragHandler(target, handler, activeClassName = "") {
    const handleMouseMove = (e) => {
        handler(target, e)
    }

    const handleMouseUp = () => {
        target.classList.remove(activeClassName)
        document.removeEventListener("mousemove", handleMouseMove)
        document.removeEventListener("mouseup", handleMouseUp)
        updateColumnSizes();
    }

    const handleMouseDown = () => {
        target.classList.add(activeClassName)
        document.addEventListener("mousemove", handleMouseMove)
        document.addEventListener("mouseup", handleMouseUp)
    }

    target.addEventListener("mousedown", handleMouseDown)
}