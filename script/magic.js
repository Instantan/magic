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
    m.socket = new WebSocket("ws://" + href);
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
        console.error('Error: ', e, 'Closing socket');
        m.socket.close();
    };
}

function handleMessage(messages) {
    const refsToRerender = new Set()
    messages.forEach(e => {
        if (e.length === 2 && typeof e[0] === 'number') {
            m.templates[e[0]] = e[1]
            makeTemplateReferenceable(e[0])
        } else if (e.length === 2 && typeof e[0] === 'string') {
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
        nanomorph(document, setDocumentClassConnectionState("connected", parseHtmlString(renderRoot())));
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
    gcRec(s, d, '0')
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
            case "magic:live":
                return magicLiveScript()
            case "magic-id":
                return v + `="${magicid}"`
        }
        if (data === undefined) {
            return ""
        }
        const toRender = data[v]
        if (toRender === undefined || toRender === null) {
            return ""
        }
        if (isRef(toRender)) {
            return renderRef(toRender)
        }
        if (isRefArray(toRender)) {
            let res = ""
            for (let i = 0; i < toRender.length; i++) {
                res += renderRef(toRender[i])
            }
            return res
        }
        return toRender
    })
}

function magicLiveScript() {
    return "<script>" + document.head.children[0].innerHTML + "</script>"
}

function execute(tpl, fn) {
    return tpl.replace(/~.+?~/g, (match) => fn(match.slice(1, -1)))
}

function parseHtmlString(markup) {
    if (markup.toLowerCase().trim().indexOf('<!doctype') === 0) {
        const doc = document.implementation.createHTMLDocument("");
        doc.documentElement.innerHTML = markup;
        return doc;
    } else if ('content' in document.createElement('template')) {
        const el = document.createElement('template');
        el.innerHTML = markup;
        return el.content;
    }
    const docfrag = document.createDocumentFragment();
    const el = document.createElement('body');
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
        nanomorph(document, hydrateTree(setDocumentClassConnectionState("connected", parseHtmlString(renderRoot()))));
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
    switch (kind) {
        case "click":
            element.onclick = createMagicEventListener(
                kind, [...baseProps], value
            )
            return
        case "focus":
            element.onfocus = createMagicEventListener(
                kind, [...baseProps], value
            )
            return
        case "change":
            element.onchange = createMagicEventListener(
                kind, [...baseProps, ...keyboardProps], value
            )
            return
        case "keydown":
            element.onkeydown = createMagicEventListener(
                kind, [...baseProps, ...keyboardProps], value
            )
            return
        case "keypress":
            element.onkeypress = createMagicEventListener(
                kind, [...baseProps, ...keyboardProps], value
            )
            return
        case "keyup":
            element.onkeyup = createMagicEventListener(
                kind, [...baseProps, ...keyboardProps], value
            )
            return
        case "submit":
            element.onsubmit = createMagicEventListener(
                kind, [...baseProps], value
            )
            return
        case "dblclick":
            element.ondblclick = createMagicEventListener(
                kind, [...baseProps], value
            )
        case "patch":
            element.onclick = liveNavigation
            return
    }
}

function cleanEvents(e) {
    e.onclick = e.onfocus = e.onchange = e.onkeydown = e.onkeypress = e.onkeyup = e.onsubmit = e.ondblclick = null
}

function createMagicEventListener(k, propsToTake, value) {
    return (e) => {
        if (m.socket) {
            m.socket.send(JSON.stringify({
                k: k,
                t: Number(getSockrefId(e.target)),
                p: Object.assign(takeFrom(e, propsToTake), { value })
            }))
        }
    }
}

function takeFrom(obj, props, value) {
    const n = {}
    let l = props.length
    while (l--) {
        const p = props[l] === "content" ? (obj.value === undefined ? obj.target.value : obj.value) : obj[props[l]]
        if (p !== undefined) {
            n[props[l]] = p
        }
    }
    return n
}

function showProgressBar() {
    if (m.topbar) {
        return
    }
    const color = m.themeColor()

    const wrapper = document.createElement("div")
    wrapper.style = `position:fixed;overflow:hidden;top:0;left:0;width:100%`

    const bg = document.createElement("div")
    bg.style = `background-color:${color};height:.2rem;width:100%;opacity:0.4`

    const inner = document.createElement("div")
    inner.style = `background-color:${color};position:absolute;bottom:0;top:0;width:50%`

    wrapper.appendChild(bg)
    wrapper.appendChild(inner)

    const timeout = setTimeout(() => {
        document.body.appendChild(wrapper)
        inner.animate([
            { left: "-50%" },
            { left: "100%" }
        ], {
            duration: 800,
            iterations: Infinity
        })
    }, 120)

    wrapper.destroy = () => {
        clearInterval(timeout)
        wrapper.remove()
    };
    m.topbar = wrapper
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

function liveNavigation(e) {
    let href = e.srcElement.attributes.getNamedItem("href").value + ""
    const path = href.startsWith("/") ? location.host + href : href
    history.pushState({}, "", href)
    connect(path)
    if (e.preventDefault) {
        e.preventDefault()
    } else if (e.stopPropagation) {
        e.stopPropagation()
    }
    return false
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
    const event = new CustomEvent(e.k, { target: e.t, detail: e.p });
    document.querySelectorAll(`[magic-id^="${e.t}"]`).forEach(e => e.dispatchEvent(event))
    window.dispatchEvent(event)
}

document.addEventListener('DOMContentLoaded', connect)