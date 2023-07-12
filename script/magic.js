import nanomorph from 'nanomorph'

window.magic = {
    templates: {},
    refs: {},
    didRenderRoot: false,
    socket: null,
    baseProps: ["metaKey", "ctrlKey", "shiftKey"],
    keyboardProps: ["key", "content"],
    themeColor: () => {
        const meta = document.querySelector('meta[name="theme-color"]')
        return meta ? meta.attributes.getNamedItem("content").value : "rgb(59,130,246)";
    }
}

const m = window.magic

function connect(href = "") {
    showProgressBar()
    setDocumentClassConnectionState("connecting")
    const previousSocket = m.socket
    if (typeof href !== "string" || href === "") {
        const url = new URL(location.host + location.pathname)
        url.searchParams = new URLSearchParams(location.search)
        href = url.toString()
    }
    m.socket = new WebSocket((window.location.protocol === "https:" ? "wss://" : "ws://") + href);
    m.socket.onopen = () => {
        if (previousSocket) {
            previousSocket.close();
        }
        m.templates = {}
        m.refs = {}
        m.didRenderRoot = false
    };
    m.socket.onmessage = (e) => {
        hideProgressBar()
        setDocumentClassConnectionState("connected")
        handleMessage(JSON.parse(e.data))
    };
    m.socket.onclose = (e) => {
        if (m.socket && e.srcElement.url !== m.socket.url) {
            return
        }
        showProgressBar()
        setDocumentClassConnectionState("disconnected")
        setTimeout(connect, 1000);
    };
    m.socket.onerror = (e) => {
        console.error("Error: ", e, "Closing socket");
        m.socket.close();
    };
}

function handleMessage(messages) {
    const refsToRerender = new Set()
    messages.forEach(e => {
        if (e.length === 2 && typeof e[0] === "number") {
            m.templates[e[0]] = e[1]
            makeTemplateReferenceable(e[0])
        } else if (e.length === 2 && typeof e[0] === "string") {
            const ref = e[0]
            assignSockref(ref, e[1])
            if (!m.didRenderRoot) {
                return
            }
            refsToRerender.add(ref)
        } else {
            receivedEvent(e)
        }
    })
    refsToRerender.forEach(updateElementsOfref)
    if (!m.didRenderRoot) {
        hydrateTree(document)
        nanomorph(document, setDocumentClassConnectionState("connected", parseHtmlString(renderRoot())), { childrenOnly: true });
        hydrateTree(document)
        m.didRenderRoot = true
    }
    gc()
}

function isRef(ref) {
    return ref && ref.length === 2 && !isNaN(ref[1])
}

function isRefArray(ref) {
    return ref && Array.isArray(ref) && ref.length > 0 && isRef(ref[0])
}

function renderRoot() {
    return renderRef(m.refs[0]['#'])
}

function gc() {
    let d = m.refs
    let s = new Set(Object.keys(d))
    gcRec(s, d, "0")
    s.forEach(e => delete d[e])
}

function gcRec(s, d, id) {
    s.delete(id)
    let o = d[id]
    if (!o) return
    for (let k in o) {
        let e = o[k]
        if (Array.isArray(e)) {
            if (isRef(e)) {
                gcRec(s, d, e[0])
            } else if (isRef(e[0])) {
                for (let i = 0; i < e.length; i++) {
                    gcRec(s, d, e[i][0])
                }
            }
        }
    }
}

function renderRef(templateref) {
    return renderTemplate(
        `${templateref[0]}:${templateref[1]}`,
        m.templates[templateref[1]],
        m.refs[templateref[0]]
    )
}

function renderTemplate(magicid, template, data) {
    return execute(template, (v) => {
        switch (v) {
            case "magic:inject":
                return magicInjects()
            case "magic-id":
                return v + `="${magicid}"`
        }
        if (data === undefined || data[v] === undefined || data[v] === null) {
            return ""
        }
        const toRender = data[v]
        if (isRef(toRender)) {
            return renderRef(toRender)
        } else if (isRefArray(toRender)) {
            let res = ""
            for (let i = 0; i < toRender.length; i++) {
                res += renderRef(toRender[i])
            }
            return res
        }
        return toRender
    })
}

function magicInjects() {
    return Array.of(...document.head.children).filter(e => e.hasAttribute("magic:inject")).map(e => {
        const n = document.createElement("div")
        n.appendChild(e)
        return n.innerHTML
    }).join('')
}

function execute(tpl, fn) {
    return tpl.replace(/~.+?~/g, (match) => fn(match.slice(1, -1)))
}

function parseHtmlString(markup) {
    if (markup.toLowerCase().trim().indexOf('<!doctype') === 0) {
        const doc = document.implementation.createHTMLDocument("");
        doc.documentElement.innerHTML = markup;
        return doc;
    } else if ('content' in createElement('template')) {
        const el = createElement('template');
        el.innerHTML = markup;
        return el.content;
    }
    const docfrag = document.createDocumentFragment();
    const el = createElement('body');
    el.innerHTML = markup;
    for (i = 0; 0 < el.childNodes.length;) {
        docfrag.appendChild(el.childNodes[i]);
    }
    return docfrag;
}

function makeTemplateReferenceable(templateid) {
    m.templates[templateid] = m.templates[templateid].replace(/<\w*(\s|>|\/>)/m, (m) => {
        let i = " ~magic-id~ "
        if (m.endsWith("/>")) {
            return m.slice(0, -1) + i + "/>"
        } else if (m.endsWith(">")) {
            return m.slice(0, -1) + i + ">"
        }
        return m + i
    })
}

function updateElementsOfref(refid) {
    if (refid === magic.refs[0]['#'][0]) {
        nanomorph(document, hydrateTree(setDocumentClassConnectionState("connected", parseHtmlString(renderRoot()))), { childrenOnly: true });
        return
    }
    document.querySelectorAll(`[magic-id^="${refid}"]`).forEach(e => {
        nanomorph(e, hydrateTree(parseHtmlString(renderRef(e.attributes.getNamedItem("magic-id").value.split(":"))).children[0]));
    })
}

function hydrateTree(tree) {
    cleanEvents(tree)
    const attrs = tree.attributes
    switch (tree.nodeName) {
        case 'TEXTAREA':
        case 'INPUT':
            tree.isSameNode = handleTextFieldValues
    }
    if (attrs) {
        for (let i = 0; i < attrs.length; i++) {
            if (attrs[i].name.startsWith("magic:")) {
                hydrateElement(tree, attrs[i])
            }
        }
    }
    if (tree.children) {
        for (let i = 0; i < tree.children.length; i++) {
            hydrateTree(tree.children[i])
        }
    }
    return tree
}

function hydrateElement(element, attribute) {
    const kind = attribute.name.slice(6)
    const value = attribute.value
    const baseProps = m.baseProps
    const keyboardProps = m.keyboardProps
    const actualEvent = "on" + kind
    switch (kind) {
        case "click":
        case "focus":
            element[actualEvent] = createMagicEventListener(
                kind, [...baseProps], value
            )
            return
        case "change":
        case "keydown":
        case "keypress":
        case "keyup":
            element[actualEvent] = createMagicEventListener(
                kind, [...baseProps, ...keyboardProps], value
            )
            return
        case "submit":
            element[actualEvent] = createMagicEventListener(
                kind, [...baseProps, "form"], value, stopPropagation
            )
            return
        case "dblclick":
            element[actualEvent] = createMagicEventListener(
                kind, [...baseProps], value
            )
            return
        case "patch":
            element[actualEvent] = liveNavigationEvent
            return
        case "static":
            element.isSameNode = (o) => {
                const s = o.attributes.getNamedItem("magic:static")
                return Boolean(s && s.value === value)
            }
            return
    }
}

function cleanEvents(e) {
    e.onclick = e.onfocus = e.onchange = e.onkeydown = e.onkeypress = e.onkeyup = e.onsubmit = e.ondblclick = null
}

function createMagicEventListener(k, propsToTake, value, after) {
    k = "m:" + k;
    return (e) => {
        if (m.socket) {
            console.log(e)
            m.socket.send(JSON.stringify({
                k: k,
                t: Number(getSockrefId(e.target)),
                p: Object.assign(takeFrom(e, propsToTake), { value })
            }))
        }
        if (typeof after === "function") {
            return after(e)
        }
    }
}

function takeFrom(obj, props, value) {
    const n = {}
    let l = props.length
    while (l--) {
        const p = props[l] === "form" ?
            Object.fromEntries((new FormData(obj.srcElement)).entries()) :
            props[l] === "content" ? (obj.value === undefined ? obj.target.value : obj.value) : obj[props[l]]
        if (p !== undefined) {
            n[props[l]] = p
        }
    }
    return n
}

function showProgressBar() {
    if (m.topbar) return

    const div = (style) => {
        const d = createElement("div")
        d.style.cssText = style
        return d
    }
    const append = (s, t) => s.appendChild(t)

    const col = m.themeColor()
    const wrp = div(`position:fixed;overflow:hidden;top:0;left:0;width:100%`)
    const inn = div(`background-color:${col};position:absolute;bottom:0;top:0;width:50%`)
    append(wrp, div(`background-color:${col};height:.2rem;width:100%;opacity:0.4`))
    append(wrp, inn)

    const timeout = setTimeout(() => {
        append(document.body, wrp)
        inn.animate([
            { left: "-50%" },
            { left: "100%" }
        ], {
            duration: 800,
            iterations: Infinity
        })
    }, 120)

    wrp.destroy = () => {
        clearInterval(timeout)
        wrp.remove()
    };
    m.topbar = wrp
}

function hideProgressBar() {
    if (m.topbar) {
        m.topbar.destroy();
        m.topbar = null;
    }
}

function assignSockref(ref, data) {
    if (m.refs[ref] === undefined) {
        m.refs[ref] = data
        return
    }
    Object.assign(m.refs[ref], data)
}

function getSockrefId(elm) {
    if (elm && elm.attributes) {
        const magicid = elm.attributes.getNamedItem("magic-id")
        if (magicid !== null) {
            return magicid.value.split(":", 1)[0]
        }
    }
    if (elm && elm.parentNode) {
        return getSockrefId(elm.parentNode)
    }
    return m.refs[0]['#'][0]
}

function handleTextFieldValues(o) {
    const n = this
    if (o.prevValue === undefined || o.prevValue !== n.value) {
        o.prevValue = n.value
    } else {
        n.value = o.value
    }
    return false
}

function liveNavigationEvent(e) {
    let href = e.srcElement.attributes.getNamedItem("href").value + ""
    liveNavigation(href)
    return stopPropagation(e)
}

function stopPropagation(e) {
    if (e.preventDefault) {
        e.preventDefault()
    } else if (e.stopPropagation) {
        e.stopPropagation()
    }
    return false
}

function liveNavigation(href) {
    const path = href.startsWith("/") ? location.host + href : href
    connect(path)
    history.pushState({}, "", href)
}

function setDocumentClassConnectionState(s, d = document) {
    if (d.children.length === 0) return
    const h = d.children[0];
    const state = "magic-" + s
    if (h.classList.contains(state)) return;
    ["connected", "connecting", "disconnected"].forEach(e => h.classList.remove("magic-" + e))
    h.classList.add(state)
    return d
}

function receivedEvent(e) {
    if (e.k.startsWith("m:")) {
        e.k = e.k.slice(2)
        switch (e.k) {
            case "navigate":
                if (m.socket) {
                    m.socket.onclose = undefined;
                }
                liveNavigation(e.p)
                return
            case "reset":
            case "click":
            case "blur":
            case "submit":
                callFnForAllElements(e.p, e.t, e.k)
                return
            case "animate":
                callFnForAllElements(e.p.ids, e.t, e.k, e.p.args.keyframes, e.p.args.options ? e.p.args.options : e.p.args.duration)
                return
            case "scrollIntoView":
                callFnForAllElements(e.p.ids, e.t, e.k, e.p.args.options ? e.p.args.options : e.p.args.alignToTop)
                return
            case "openFullscreen":
                openFullscreen()
                return
            case "closeFullscreen":
                closeFullscreen()
                return
            case "disconnect":
                if (m.socket) {
                    m.socket.onclose = undefined;
                    m.socket.close();
                    setDocumentClassConnectionState("disconnected")
                }
                return
            case "updateUrl":
                history.pushState({}, "", e.p)
            case "refreshFile":
                callFnForAllElements(e.p, e.t, refreshFile)
                return
        }
    }
    const event = new CustomEvent(e.k, { target: e.t, detail: e.p });
    document.querySelectorAll(`[magic-id^="${e.t}"]`).forEach(e => e.dispatchEvent(event))
    window.dispatchEvent(event)
}

function createElement(tag) {
    return document.createElement(tag)
}

function callFnForAllElements(ids = [], socket = null, fn = "", ...args) {
    const call = (elm) => {
        try {
            if (typeof fn === 'function') {
                fn(elm)
            } else if (elm && elm[fn]) {
                elm[fn](...args)
            }
        } catch (ex) {
            console.error(ex)
        }
    }
    ids.forEach(id => call(document.getElementById(id)))
    if (ids.length === 0) {
        document.querySelectorAll(`[magic-id^="${e.t}"]`).forEach(call)
    }
}

function openFullscreen() {
    if (document.requestFullscreen) {
        document.requestFullscreen();
    } else if (elem.webkitRequestFullscreen) { /* Safari */
        document.webkitRequestFullscreen();
    } else if (elem.msRequestFullscreen) { /* IE11 */
        document.msRequestFullscreen();
    }
}

function closeFullscreen() {
    if (document.exitFullscreen) {
        document.exitFullscreen();
    } else if (document.webkitExitFullscreen) { /* Safari */
        document.webkitExitFullscreen();
    } else if (document.msExitFullscreen) { /* IE11 */
        document.msExitFullscreen();
    }
}

function refreshFile(elm) {
    if (elm.src) {
        let breaker = elm.src.includes("?") && !elm.src.includes("?t=") ? "&t=" : "?t="
        let src = elm.src;
        if (src.indexOf(breaker) != -1) {
            src = src.split(breaker)[0];
        }
        elm.src = src + breaker + Date.now();
    }
}

document.addEventListener('DOMContentLoaded', connect)