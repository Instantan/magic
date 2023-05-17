/* @refresh reload */
import './../assets/css/index.css'

document.getElementById("open-menu").addEventListener("click", () => {
    document.getElementById("menu").style.display = "block"
})
document.getElementById("close-menu").addEventListener("click", () => {
    document.getElementById("menu").style.display = "none"
})

const hashchange = () => {
    document.body.classList.forEach(l => {
        if (l.startsWith("hash-")) {
            document.body.classList.remove(l)
        }
    })
    if (window.location.hash.length === 0) {
        return
    }
    document.body.classList.add("hash-" + window.location.hash.slice(1))
}
hashchange()
window.addEventListener('hashchange', hashchange)