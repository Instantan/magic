import nanomorph from 'nanomorph'

window.magic = {
    templates: {},
    socketrefs: {},
    didRenderRoot: false,
    socket: null
}

function connect() {
    const ws_params = new URLSearchParams(location.search);
    ws_params.append("ws", "0");
    window.magic.socket = new WebSocket("ws://" + location.host + location.pathname + "?" + ws_params);
    window.magic.socket.onopen = function() {
        window.magic.didRenderRoot = false
    };
    window.magic.socket.onmessage = function(e) {
        handleMessage(JSON.parse(e.data))
    };
    window.magic.socket.onclose = function(e) {
        console.log('Socket is closed. Reconnect will be attempted in 1 second.', e.reason);
        setTimeout(function() {
            connect();
        }, 1000);
    };
    window.magic.socket.onerror = function(err) {
        console.error('Socket encountered error: ', err, 'Closing socket');
        magic.socket.close();
    };
}
connect()

function handleMessage(message) {
    const socketrefsToRerender = new Set()
    message.forEach(element => {
        console.log(element)
        if (element.length === 2 && typeof element[0] === 'number') {
            window.magic.templates[element[0]] = element[1]
            makeTemplateReferenceable(element[0])
        } else if (element.length === 2 && isSocketId(element[0])) {
            const ref = element[0]
            if (window.magic.socketrefs[ref] === undefined) {
                window.magic.socketrefs[ref] = element[1]
            }
            Object.assign(window.magic.socketrefs[ref], element[1])
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
        window.magic.didRenderRoot = true
    }
}

function isSocketId(socketid) {
    return typeof socketid === 'string'
}

function isTemplateRef(ref) {
    return ref.length === 2 && typeof ref[0] === 'string'
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
        if (isTemplateRef(toRender)) {
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
        return
    }
    document.querySelectorAll(`[magic-id^="${socketrefid}"]`).forEach(elm => {
        const newElm = parseHtmlString(renderTemplateRef(elm.attributes.getNamedItem("magic-id").value.split(":")))
        morph(elm, newElm.children[0])
    })
}