/* @refresh reload */
import './../index.css'

const twlightmode = 'tw-lightmode'

function handleLightSwitch() {
    const isLight = document.documentElement.classList.contains('light')
    const setTo = isLight ? 'dark' : 'light'
    localStorage.setItem(twlightmode, setTo)
    document.documentElement.classList.replace(setTo === 'light' ? 'dark' : 'light', setTo)
}

function handleNavbarSwitch() {
    const nav = document.getElementsByTagName('nav')[0]
    nav.style.display = nav.style.display === 'none' ? '' : 'none'
}

function handleEyecatcherSwitch() {
    document.getElementById('eyecatcher').classList.toggle('full')
}

function startEyecatcherAnimations() {
    // Some random colors
    const colors = ["#2563eb", "#52525b", "#34d399", "#fb923c", "#db2777"];

    const numBalls = 100;
    const balls = [];

    const target = document.getElementById("eyecatcher")

    for (let i = 0; i < numBalls; i++) {
        let ball = document.createElement("div");
        ball.classList.add("point");
        ball.style.background = colors[Math.floor(Math.random() * colors.length)];
        ball.style.left = `${Math.floor(Math.random() * 100)}vw`;
        ball.style.top = `${Math.floor(Math.random() * 100)}vh`;
        ball.style.transform = `scale(${Math.random()})`;
        ball.style.width = `${Math.random()}em`;
        ball.style.height = ball.style.width;

        balls.push(ball);
        target.append(ball);
    }

    // Keyframes
    balls.forEach((el, i, ra) => {
        let to = {
            x: Math.random() * (i % 2 === 0 ? -11 : 11),
            y: Math.random() * 12
        };

        let anim = el.animate(
            [
                { transform: "translate(0, 0)" },
                { transform: `translate(${to.x}rem, ${to.y}rem)` }
            ], {
                duration: (Math.random() + 1) * 500, // random duration
                direction: "alternate",
                fill: "both",
                iterations: Infinity,
                easing: "ease-in-out"
            }
        );
    });
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
        if ((e.metaKey || e.ctrlKey) && e.key === 'b') {
            handleNavbarSwitch()
            handled = true
        }
        if ((e.metaKey || e.ctrlKey) && e.key === 'j') {
            handleEyecatcherSwitch()
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
    addShortcutEventListener()
    lightSwitchEventListener()
    startEyecatcherAnimations()
});