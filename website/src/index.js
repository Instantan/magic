/* @refresh reload */
import './../assets/css/index.css'

document.getElementById("open-menu").addEventListener("click", () => {
    document.getElementById("menu").style.display = "block"
})
document.getElementById("close-menu").addEventListener("click", () => {
    document.getElementById("menu").style.display = "none"
})