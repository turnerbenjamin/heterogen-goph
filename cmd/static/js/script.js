//Helpers
function isSuccess(statusCode){
    return statusCode >= 200 && statusCode < 400
}

//ANIMATE TOAST
let toastTimeout
const toastContainer = document.getElementById("toast-container")
const toastEl = document.getElementById("toast")
const toastmessage = document.getElementById("toast-message")
const toastIcon = document.getElementById("toast-icon")

function hideToast(){
    if(!toastTimeout) return;
    toastEl.classList.add("hide")
    toastEl.classList.remove("show")
    clearTimeout(toastTimeout)
    toastTimeout = undefined;
}

function updateToast(){
    if(toastTimeout) hideToast()
    const toastmessage = document.getElementById("toast-message")
    if(!toastmessage?.innerText) return

    const containerTop = toastContainer.getBoundingClientRect().top
    const containerPadding = containerTop < 5 ? 5 - containerTop : 0;
    toastContainer.style.paddingTop = containerPadding + "px"; 

    toastEl.classList.add("show")
    toastEl.classList.remove("hide")
    toastTimeout = setTimeout(hideToast, 7600)
}

function setSuccessToast(msg, status = 200){
    if(!isSuccess(status)) return
    toastmessage.innerHTML = msg
    toastIcon.innerText = "✓"
    toastEl.classList.add("show", "success")
    updateToast()
}

function setErrorToast(msg){
    console.log(msg)
    toastmessage.innerHTML = msg
    toastIcon.innerText = "✖"
    toastEl.classList.add("show", "error")
    updateToast()
}

// Toggle Sign-In Form
const signInFormEl = document.getElementById("sign-in-form-container")
function showSignInForm(){
    signInFormEl.classList.remove("hidden")
}

function hideSignInForm(){
    signInFormEl.classList.add("hidden")
}


//PERSIST STATE AFTER REDIRECT/Refresh
function setState(key, value){
    localStorage.setItem(key, value)
    hideToast()
}
function removeState(key){  
    localStorage.removeItem(key)
}

//LOG-IN STATE
const isLoggingInKey = "is-logging-In"
const setIsLoggingIn = () => setState(isLoggingInKey, "true")
const clearIsLoggingIn = (status) => {
    if(status >=200 && status <400) return
    removeState(isLoggingInKey)
}

const didLogIn = localStorage.getItem(isLoggingInKey)
if(didLogIn){
    clearIsLoggingIn();
    setSuccessToast("You have logged-in successfully", 200)
} 

//LOG OUT STATE
const isLoggingOutKey = "is-logging-out"
function setIsLoggingOut(){setState(isLoggingOutKey, "true")}
function clearIsLoggingOut(){removeState(isLoggingOutKey)}

const didLogOut = localStorage.getItem(isLoggingOutKey)
if(didLogOut){
    clearIsLoggingOut();
    setSuccessToast("You have logged-out successfully", 200)
} 

//REGISTRATION STATE
const isRegisteringKey = "is-registering"
function setIsRegistering(){setState(isRegisteringKey, "true")}
function clearIsRegistering(){removeState(isRegisteringKey)}

const didRegister = localStorage.getItem(isRegisteringKey)
if(didRegister){
    clearIsRegistering();
    setSuccessToast(`You have created an account successfully.`, 200)
    showSignInForm()
} 








//GLOBAL ERROR HANDLING
function handleErrors(e){  
    if (isSuccess(e.detail.xhr.status)) return
    const errorMessageContainerEl = e.target.getElementsByClassName("error-message-container") 
    const errorMarkup = e.detail.xhr.response
    console.l    
    if(errorMessageContainerEl.length > 0){
        errorMessageContainerEl[0].innerHTML = errorMarkup
    }else{
        setErrorToast(errorMarkup); 
    }
    clearIsLoggingOut();
}
document.addEventListener("htmx:afterRequest", handleErrors)








