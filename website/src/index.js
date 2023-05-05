/* @refresh reload */
import './../index.css'

const twlightmode = 'tw-lightmode'

function handleLightSwitch() {
    const isLight = document.documentElement.classList.contains('light')
    const setTo = isLight ? 'dark' : 'light'
    localStorage.setItem(twlightmode, setTo)
    document.documentElement.classList.replace(setTo === 'light' ? 'dark' : 'light', setTo)
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
});