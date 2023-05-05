import nanomorph from 'nanomorph'

window.magic = {
    templates: new Map(),
    socketrefs: new Map(),
    socketrefs_refs: new Map(),
    didRenderRoot: false,
    socket: null,
    baseProps: ["metaKey", "ctrlKey", "shiftKey"],
    themeColor: () => {
        const meta = document.querySelector('meta[name="theme-color"]')
        return meta ? meta.attributes.getNamedItem("content").value : "rgb(59,130,246)";
    }
}

function connect() {
    showProgressBar()
    const ws_params = new URLSearchParams(location.search);
    ws_params.append("ws", "0");
    window.magic.socket = new WebSocket("ws://" + location.host + location.pathname + "?" + ws_params);
    window.magic.socket.onopen = () => {
        const m = window.magic
        m.templates = {}
        m.socketrefs = {}
        m.socketrefs_refs = {}
        m.didRenderRoot = false
    };
    window.magic.socket.onmessage = (e) => {
        hideProgressBar()
        handleMessage(JSON.parse(e.data))
    };
    window.magic.socket.onclose = (e) => {
        showProgressBar()
        setTimeout(connect, 1000);
    };
    window.magic.socket.onerror = (e) => {
        console.error('Socket encountered error: ', e, 'Closing socket');
        magic.socket.close();
    };
}

function handleMessage(message) {
    const socketrefsToRerender = new Set()
    message.forEach(element => {
        if (element.length === 2 && typeof element[0] === 'number') {
            window.magic.templates[element[0]] = element[1]
            makeTemplateReferenceable(element[0])
        } else if (element.length === 2 && isSocketId(element[0])) {
            const ref = element[0]
            assignSockref(ref, element[1])
            if (!window.magic.didRenderRoot) {
                return
            }
            socketrefsToRerender.add(ref)
        }
    })
    socketrefsToRerender.forEach(ref => {
        updateElementsOfSocketref(ref)
    })
    if (!window.magic.didRenderRoot) {
        nanomorph(document, parseHtmlString(renderRoot()));
        hydrateTree(document)
        window.magic.didRenderRoot = true
    }
}

function isSocketId(socketid) {
    return typeof socketid === 'string'
}

function isRef(ref) {
    return ref && ref.length === 2 && typeof ref[0] === 'string'
}

function renderRoot() {
    return renderTemplateRef(window.magic.socketrefs[0]['#'])
}

function renderTemplateRef(templateref) {
    return renderTemplate(
        `${templateref[0]}:${templateref[1]}`,
        window.magic.templates[templateref[1]],
        window.magic.socketrefs[templateref[0]]
    )
}

function renderTemplate(magicid, template, data) {
    return execute(template, (v) => {
        switch (v) {
            case "magic:live":
                return magicLiveScript()
            case "magic:id":
                return `magic-id="${magicid}"`
        }
        if (data === undefined) {
            return ""
        }
        const toRender = data[v]
        if (toRender === undefined || toRender === null) {
            return ""
        }
        if (isRef(toRender)) {
            return renderTemplateRef(toRender)
        }
        return toRender
    })
}

function magicLiveScript() {
    return "<script>" + document.head.children[0].innerHTML + "</script>"
}

function execute(tpl, fn) {
    return tpl.replace(/~.+?~/g, (match, contents, offset, input_string) => fn(match.slice(1, -1)))
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
    window.magic.templates[templateid] = window.magic.templates[templateid].replace(/<\w*(\s|>|\/>)/m, (m) => {
        if (m.endsWith("/>")) {
            return m.slice(0, -1) + " ~magic:id~/>"
        } else if (m.endsWith(">")) {
            return m.slice(0, -1) + " ~magic:id~>"
        }
        return m + "~magic:id~ "
    })
}

function updateElementsOfSocketref(socketrefid) {
    if (socketrefid === magic.socketrefs[0]['#'][0]) {
        nanomorph(document, parseHtmlString(renderRoot()));
        hydrateTree(document)
        console.debug("[RENDERED]", document)
        return
    }
    document.querySelectorAll(`[magic-id^="${socketrefid}"]`).forEach(elm => {
        const newElm = parseHtmlString(renderTemplateRef(elm.attributes.getNamedItem("magic-id").value.split(":")))
        nanomorph(elm, newElm.children[0]);
        hydrateTree(elm);
        console.debug("[RENDERED]", elm);
    })
}

function hydrateTree(tree) {
    cleanEvents(tree)
    const attrs = tree.attributes
    if (attrs) {
        for (let i = 0; i < attrs.length; i++) {
            if (attrs[i].name.startsWith("magic")) {
                hydrateElement(tree, attrs[i])
            }
        }
    }
    if (tree.children) {
        for (let i = 0; i < tree.children.length; i++) {
            hydrateTree(tree.children[i])
        }
    }
}

function hydrateElement(element, attribute) {
    const kind = attribute.name.slice(6)
    const value = attribute.value
    const baseProps = window.magic.baseProps
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
                kind, [...baseProps], value
            )
            return
        case "keydown":
            element.onkeydown = createMagicEventListener(
                kind, [...baseProps], value
            )
            return
        case "keypress":
            element.onkeypress = createMagicEventListener(
                kind, [...baseProps], value
            )
            return
        case "keyup":
            element.onkeyup = createMagicEventListener(
                kind, [...baseProps], value
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
            return
    }
}

function cleanEvents(e) {
    e.onclick = null
    e.onfocus = null
    e.onchange = null
    e.onkeydown = null
    e.onkeypress = null
    e.onkeyup = null
    e.onsubmit = null
    e.ondblclick = null
}

function createMagicEventListener(kind, propsToTake, value) {
    return (e) => {
        if (window.magic.socket) {
            const target = Number(getSockrefId(e.target))
            const payload = Object.assign(takeFrom(e, propsToTake), { value })
            window.magic.socket.send(JSON.stringify({
                kind,
                target,
                payload,
            }))
        }
    }
}

function takeFrom(obj, props, value) {
    const n = {}
    let l = props.length
    while (l--) {
        const p = obj[props[l]]
        if (p !== undefined) {
            n[props[l]] = p
        }
    }
    return n
}

function showProgressBar() {
    if (window.magic.topbar) {
        return
    }

    const color = window.magic.themeColor()

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
    window.magic.topbar = wrapper
}

function hideProgressBar() {
    if (window.magic.topbar) {
        window.magic.topbar.destroy();
        window.magic.topbar = null;
    }
}

function assignSockref(ref, data) {
    const newFields = Object.keys(data)
    let nfl = newFields.length;
    while (nfl--) {
        const v = data[newFields[nfl]]
        if (isRef(v)) {
            socketrefTrack(v[0], +1)
        }
    }
    if (window.magic.socketrefs[ref] === undefined) {
        window.magic.socketrefs[ref] = data
        return
    }
    nfl = newFields.length;
    while (nfl--) {
        const v = window.magic.socketrefs[ref][newFields[nfl]]
        if (isRef(v)) {
            socketrefTrack(v[0], -1)
        }
    }
    Object.assign(window.magic.socketrefs[ref], data)
}

function socketrefTrack(ref, action) {
    if (window.magic.socketrefs_refs[ref] !== undefined) {
        window.magic.socketrefs_refs[ref] += action
    } else {
        window.magic.socketrefs_refs[ref] = action
    }
    if (action < 0 && window.magic.socketrefs_refs[ref] < 1) {
        delete window.magic.socketrefs[ref];
        delete window.magic.socketrefs_refs[ref];
    }
}

function getSockrefId(elm) {
    const magicid = elm.attributes.getNamedItem("magic-id")
    if (magicid !== null) {
        return magicid.value.split(":", 1)[0]
    }
    if (elm.parentNode) {
        return getSockrefId(elm)
    }
    return window.magic.socketrefs[0]['#'][0]
}

document.addEventListener('DOMContentLoaded', connect)