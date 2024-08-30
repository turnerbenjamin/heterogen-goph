//TOAST ANIMATION

let toastTimeout

function hideToast(){
    const toastEl = document.getElementById("toast")
    toastEl.classList.remove("show")
    toastEl.classList.add("hide")
    clearTimeout(toastTimeout)
}

function updateToast(){
    const toastEl = document.getElementById("toast")
    const toastmessage = document.getElementById("toast-message")
    console.log(toastmessage.innerText)
    if(!toastmessage.innerHTML) return
    toastEl.classList.add("show")
    toastEl.classList.remove("hide")
    toastTimeout = setTimeout(hideToast, 7600)
}




