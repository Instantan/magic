/* @refresh reload */
import './../index.css'

const twlightmode = 'tw-lightmode'

function handleLightSwitch() {
    const isLight = document.documentElement.classList.contains('light')
    const setTo = isLight ? 'dark' : 'light'
    localStorage.setItem(twlightmode, setTo)
    document.documentElement.classList.replace(setTo === 'light' ? 'dark' : 'light', setTo)
}

function startEyecatcherAnimations() {
    const shape1 = document.getElementById("shape1")
    const shape2 = document.getElementById("shape2")
    const shape3 = document.getElementById("shape3")
    window.onmousemove = function(e) {
        shape1.animate({
            left: `${e.clientX - 100}px`,
            top: `${e.clientY}px`
        }, {
            duration: 6000,
            fill: "forwards"
        })
        shape2.animate({
            left: `${e.clientX + 200}px`,
            top: `${e.clientY + 200}px`
        }, {
            duration: 4000,
            fill: "forwards"
        })
        shape3.animate({
            left: `${e.clientX}px`,
            top: `${e.clientY - 200}px`
        }, {
            duration: 8000,
            fill: "forwards"
        })
    }
}


function lightSwitchEventListener() {
    Array.from(document.getElementsByClassName("lightswitch")).forEach(function(e) {
        e.addEventListener("click", function(event) {
            handleLightSwitch()
        })
    })
}

function addShortcutEventListener() {
    window.addEventListener('keydown', (e) => {
        let handled = false
        if ((e.metaKey || e.ctrlKey) && e.key === 'l') {
            handleLightSwitch()
            handled = true
        }
        if (handled) {
            if (e && e.stopPropagation) {
                e.stopPropagation();
            } else if (e && window.event) {
                window.event.cancelBubble = true;
            }
            e.preventDefault();
        }
    })
}

document.addEventListener("DOMContentLoaded", function(event) {
    lightSwitchEventListener()
    addShortcutEventListener()
    startEyecatcherAnimations()
});