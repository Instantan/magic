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
    window.magic.socket.onopen = function() {
        window.magic.templates = {}
        window.magic.socketrefs = {}
        window.magic.socketrefs_refs = {}
        window.magic.didRenderRoot = false
    };
    window.magic.socket.onmessage = function(e) {
        hideProgressBar()
        handleMessage(JSON.parse(e.data))
    };
    window.magic.socket.onclose = function(e) {
        showProgressBar()
        setTimeout(connect, 1000);
    };
    window.magic.socket.onerror = function(err) {
        console.error('Socket encountered error: ', err, 'Closing socket');
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
        morph(document, parseHtmlString(renderRoot()))
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
    const ref = magic.socketrefs[0]['#']
    return renderTemplateRef(ref)
}

function renderTemplateRef(templateref) {
    return renderTemplate(
        `${templateref[0]}:${templateref[1]}`,
        window.magic.templates[templateref[1]],
        window.magic.socketrefs[templateref[0]]
    )
}

function renderTemplate(magicid, template, data) {
    const res = execute(template, (v) => {
        if (v === "magic:live") {
            return magicLiveScript()
        } else if (v === "magic:id") {
            return `magic-id="${magicid}"`
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
    return res
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
    } else {
        const docfrag = document.createDocumentFragment();
        const el = document.createElement('body');
        el.innerHTML = markup;
        for (i = 0; 0 < el.childNodes.length;) {
            docfrag.appendChild(el.childNodes[i]);
        }
        return docfrag;
    }
}

function morph(oldTree, newTree) {
    nanomorph(oldTree, newTree)
}

function makeTemplateReferenceable(templateid) {
    window.magic.templates[templateid] = window.magic.templates[templateid].replace(/<\w*(\s|>|\/>)/m, (match, contents, offset, input_string) => {
        if (match.endsWith("/>")) {
            return match.slice(0, -1) + " ~magic:id~/>"
        } else if (match.endsWith(">")) {
            return match.slice(0, -1) + " ~magic:id~>"
        } else {
            return match + "~magic:id~ "
        }
        return match
    })
}

function updateElementsOfSocketref(socketrefid) {
    if (socketrefid === magic.socketrefs[0]['#'][0]) {
        morph(document, parseHtmlString(renderRoot()))
        hydrateTree(document)
        console.debug("[RENDERED]", document)
        return
    }
    document.querySelectorAll(`[magic-id^="${socketrefid}"]`).forEach(elm => {
        const newElm = parseHtmlString(renderTemplateRef(elm.attributes.getNamedItem("magic-id").value.split(":")))
        morph(elm, newElm.children[0])
        hydrateTree(elm)
        console.debug("[RENDERED]", elm)
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

function cleanEvents(element) {
    element.onclick = null
    element.onfocus = null
    element.onchange = null
    element.onkeydown = null
    element.onkeypress = null
    element.onkeyup = null
    element.onsubmit = null
    element.ondblclick = null
}

function createMagicEventListener(kind, propsToTake, value) {
    return (e) => {
        const payload = Object.assign(takeFrom(e, propsToTake), { value })
        if (window.magic.socket) {
            // somehow we need to get the event target
            window.magic.socket.send(JSON.stringify({
                kind,
                payload
            }))
        }
    }
}

function takeFrom(obj, props, value) {
    const n = {}
    for (let i = 0; i < props.length; i++) {
        const p = obj[props[i]]
        if (p !== undefined) {
            n[props[i]] = p
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
    ref = Number(ref)
    const newFields = Object.keys(data)
    let nfl = newFields.length;
    // incr action
    while (nfl--) {
        const v = data[newFields[nfl]]
        if (isRef(v)) {
            socketrefTrack(v[0], +1)
        }
    }

    if (window.magic.socketrefs[ref] === undefined) {
        window.magic.socketrefs[ref] = data
    } else {
        nfl = newFields.length;
        // decr action
        while (nfl--) {
            const v = window.magic.socketrefs[ref][newFields[nfl]]
            if (isRef(v)) {
                socketrefTrack(v[0], -1)
            }
        }
        Object.assign(window.magic.socketrefs[ref], data)
    }
}

function socketrefTrack(ref, action) {
    ref = Number(ref)
    if (window.magic.socketrefs_refs[ref] !== undefined) {
        window.magic.socketrefs_refs[ref] += action
    } else {
        window.magic.socketrefs_refs[ref] = action
    }
    if (action < 0 && window.magic.socketrefs_refs[ref] < 1) {
        console.info(delete window.magic.socketrefs[ref]);
        console.info(delete window.magic.socketrefs_refs[ref]);
        console.log("Found orphan and removed it:", ref, Object.keys(window.magic.socketrefs_refs).length, Object.keys(window.magic.socketrefs).length)
    }
}

document.addEventListener('DOMContentLoaded', connect)