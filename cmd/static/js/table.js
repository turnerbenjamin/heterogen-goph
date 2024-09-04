
let headers;
let tableConfig

//Initial Load
document.addEventListener('DOMContentLoaded', () => {
    htmx.trigger("#table-container", "refresh-table")
});

function handleTableRefresh(){
    headers = document.querySelectorAll("th")
    addTableHandlers()
    config = getTableConfig();
    if(!config) {
        config = initialiseTableConfig()
    }
    applyConfig()   
}

function addTableHandlers(){ 
    const table = document.querySelector("table");
    if(!table) return
    const tableResizers = table.querySelectorAll(".resizer")
    
    const tableHeight = table.getBoundingClientRect().height;
    for (const tableResizer of tableResizers) {
        tableResizer.style.height = tableHeight + "px"
        handleDrag(tableResizer, handleResize)
    }
}

function applyConfig(){
    for(const header of headers){
        const headerName = header.innerText;
        const headerConfig = config.columnConfig[headerName]
        
        header.style.width = headerConfig.width + "px"
    }
    headers.forEach(header=>{
        headerName = header.innerText;
        headerConfig = {
            width: header.getBoundingClientRect().width
        }
    })
}

function updateColumnSizes(){

    headers.forEach(header=>{
        const headerName = header.innerText;
        const headerConfig = config.columnConfig[headerName]
        const width = header.getBoundingClientRect().width;
        headerConfig.width = width
    })
    localStorage.setItem("userTableConfig", JSON.stringify(config))
}


function getTableConfig(){
    const config = localStorage.getItem("userTableConfig")
    if(!config) return ""
    return JSON.parse(config)
}





function initialiseTableConfig(){
    const tableConfig = {
        columns: [],
        columnConfig: {}
    }

    headers.forEach(header=>{
        headerName = header.innerText;
        headerConfig = {
            width: header.getBoundingClientRect().width
        }
        tableConfig.columns.push(headerName);
        tableConfig.columnConfig[headerName] = headerConfig
    })
    localStorage.setItem("userTableConfig", JSON.stringify(tableConfig))
    return tableConfig
}



function handleResize(target, e){
    const resizer = target
    const header = resizer.closest("th")
    const { width, right } = header.getBoundingClientRect();
    const dx = e.clientX- right;
    const newWidth = width + dx    
    header.style.width = newWidth + 'px'; 
}


function handleDrag(target, handler){
    const handleMouseMove = (e)=>{
        handler(target, e)  
       }
       
       const handleMouseUp = ()=>{ 
        target.classList.remove("dragging")
           document.removeEventListener("mousemove", handleMouseMove)
           document.removeEventListener("mouseup", handleMouseUp)
           updateColumnSizes();
           htmx.trigger("#table-container", "refresh-table", {test:15})
       }

       const handleMouseDown = ()=>{
        target.classList.add("dragging")
        document.addEventListener("mousemove", handleMouseMove)
        document.addEventListener("mouseup", handleMouseUp)
       }

       target.addEventListener("mousedown", handleMouseDown)      
}