(function(global, factory) { typeof exports === "object" && typeof module !== "undefined" ? module.exports = factory() : typeof define === "function" && define.amd ? define(factory) : (global = global || self, global.morphdom = factory()) })(this, function() {
    "use strict";
    var DOCUMENT_FRAGMENT_NODE = 11;

    function morphAttrs(fromNode, toNode) {
        var toNodeAttrs = toNode.attributes;
        var attr;
        var attrName;
        var attrNamespaceURI;
        var attrValue;
        var fromValue;
        if (toNode.nodeType === DOCUMENT_FRAGMENT_NODE || fromNode.nodeType === DOCUMENT_FRAGMENT_NODE) { return }
        for (var i = toNodeAttrs.length - 1; i >= 0; i--) {
            attr = toNodeAttrs[i];
            attrName = attr.name;
            attrNamespaceURI = attr.namespaceURI;
            attrValue = attr.value;
            if (attrNamespaceURI) {
                attrName = attr.localName || attrName;
                fromValue = fromNode.getAttributeNS(attrNamespaceURI, attrName);
                if (fromValue !== attrValue) {
                    if (attr.prefix === "xmlns") { attrName = attr.name }
                    fromNode.setAttributeNS(attrNamespaceURI, attrName, attrValue)
                }
            } else { fromValue = fromNode.getAttribute(attrName); if (fromValue !== attrValue) { fromNode.setAttribute(attrName, attrValue) } }
        }
        var fromNodeAttrs = fromNode.attributes;
        for (var d = fromNodeAttrs.length - 1; d >= 0; d--) {
            attr = fromNodeAttrs[d];
            attrName = attr.name;
            attrNamespaceURI = attr.namespaceURI;
            if (attrNamespaceURI) { attrName = attr.localName || attrName; if (!toNode.hasAttributeNS(attrNamespaceURI, attrName)) { fromNode.removeAttributeNS(attrNamespaceURI, attrName) } } else { if (!toNode.hasAttribute(attrName)) { fromNode.removeAttribute(attrName) } }
        }
    }
    var range;
    var NS_XHTML = "http://www.w3.org/1999/xhtml";
    var doc = typeof document === "undefined" ? undefined : document;
    var HAS_TEMPLATE_SUPPORT = !!doc && "content" in doc.createElement("template");
    var HAS_RANGE_SUPPORT = !!doc && doc.createRange && "createContextualFragment" in doc.createRange();

    function createFragmentFromTemplate(str) {
        var template = doc.createElement("template");
        template.innerHTML = str;
        return template.content.childNodes[0]
    }

    function createFragmentFromRange(str) {
        if (!range) {
            range = doc.createRange();
            range.selectNode(doc.body)
        }
        var fragment = range.createContextualFragment(str);
        return fragment.childNodes[0]
    }

    function createFragmentFromWrap(str) {
        var fragment = doc.createElement("body");
        fragment.innerHTML = str;
        return fragment.childNodes[0]
    }

    function toElement(str) { str = str.trim(); if (HAS_TEMPLATE_SUPPORT) { return createFragmentFromTemplate(str) } else if (HAS_RANGE_SUPPORT) { return createFragmentFromRange(str) } return createFragmentFromWrap(str) }

    function compareNodeNames(fromEl, toEl) {
        var fromNodeName = fromEl.nodeName;
        var toNodeName = toEl.nodeName;
        var fromCodeStart, toCodeStart;
        if (fromNodeName === toNodeName) { return true }
        fromCodeStart = fromNodeName.charCodeAt(0);
        toCodeStart = toNodeName.charCodeAt(0);
        if (fromCodeStart <= 90 && toCodeStart >= 97) { return fromNodeName === toNodeName.toUpperCase() } else if (toCodeStart <= 90 && fromCodeStart >= 97) { return toNodeName === fromNodeName.toUpperCase() } else { return false }
    }

    function createElementNS(name, namespaceURI) { return !namespaceURI || namespaceURI === NS_XHTML ? doc.createElement(name) : doc.createElementNS(namespaceURI, name) }

    function moveChildren(fromEl, toEl) {
        var curChild = fromEl.firstChild;
        while (curChild) {
            var nextChild = curChild.nextSibling;
            toEl.appendChild(curChild);
            curChild = nextChild
        }
        return toEl
    }

    function syncBooleanAttrProp(fromEl, toEl, name) { if (fromEl[name] !== toEl[name]) { fromEl[name] = toEl[name]; if (fromEl[name]) { fromEl.setAttribute(name, "") } else { fromEl.removeAttribute(name) } } }
    var specialElHandlers = {
        OPTION: function(fromEl, toEl) {
            var parentNode = fromEl.parentNode;
            if (parentNode) {
                var parentName = parentNode.nodeName.toUpperCase();
                if (parentName === "OPTGROUP") {
                    parentNode = parentNode.parentNode;
                    parentName = parentNode && parentNode.nodeName.toUpperCase()
                }
                if (parentName === "SELECT" && !parentNode.hasAttribute("multiple")) {
                    if (fromEl.hasAttribute("selected") && !toEl.selected) {
                        fromEl.setAttribute("selected", "selected");
                        fromEl.removeAttribute("selected")
                    }
                    parentNode.selectedIndex = -1
                }
            }
            syncBooleanAttrProp(fromEl, toEl, "selected")
        },
        INPUT: function(fromEl, toEl) {
            syncBooleanAttrProp(fromEl, toEl, "checked");
            syncBooleanAttrProp(fromEl, toEl, "disabled");
            if (fromEl.value !== toEl.value) { fromEl.value = toEl.value }
            if (!toEl.hasAttribute("value")) { fromEl.removeAttribute("value") }
        },
        TEXTAREA: function(fromEl, toEl) {
            var newValue = toEl.value;
            if (fromEl.value !== newValue) { fromEl.value = newValue }
            var firstChild = fromEl.firstChild;
            if (firstChild) {
                var oldValue = firstChild.nodeValue;
                if (oldValue == newValue || !newValue && oldValue == fromEl.placeholder) { return }
                firstChild.nodeValue = newValue
            }
        },
        SELECT: function(fromEl, toEl) {
            if (!toEl.hasAttribute("multiple")) {
                var selectedIndex = -1;
                var i = 0;
                var curChild = fromEl.firstChild;
                var optgroup;
                var nodeName;
                while (curChild) {
                    nodeName = curChild.nodeName && curChild.nodeName.toUpperCase();
                    if (nodeName === "OPTGROUP") {
                        optgroup = curChild;
                        curChild = optgroup.firstChild
                    } else {
                        if (nodeName === "OPTION") {
                            if (curChild.hasAttribute("selected")) { selectedIndex = i; break }
                            i++
                        }
                        curChild = curChild.nextSibling;
                        if (!curChild && optgroup) {
                            curChild = optgroup.nextSibling;
                            optgroup = null
                        }
                    }
                }
                fromEl.selectedIndex = selectedIndex
            }
        }
    };
    var ELEMENT_NODE = 1;
    var DOCUMENT_FRAGMENT_NODE$1 = 11;
    var TEXT_NODE = 3;
    var COMMENT_NODE = 8;

    function noop() {}

    function defaultGetNodeKey(node) { if (node) { return node.getAttribute && node.getAttribute("id") || node.id } }

    function morphdomFactory(morphAttrs) {
        return function morphdom(fromNode, toNode, options) {
            if (!options) { options = {} }
            if (typeof toNode === "string") {
                if (fromNode.nodeName === "#document" || fromNode.nodeName === "HTML" || fromNode.nodeName === "BODY") {
                    var toNodeHtml = toNode;
                    toNode = doc.createElement("html");
                    toNode.innerHTML = toNodeHtml
                } else { toNode = toElement(toNode) }
            }
            var getNodeKey = options.getNodeKey || defaultGetNodeKey;
            var onBeforeNodeAdded = options.onBeforeNodeAdded || noop;
            var onNodeAdded = options.onNodeAdded || noop;
            var onBeforeElUpdated = options.onBeforeElUpdated || noop;
            var onElUpdated = options.onElUpdated || noop;
            var onBeforeNodeDiscarded = options.onBeforeNodeDiscarded || noop;
            var onNodeDiscarded = options.onNodeDiscarded || noop;
            var onBeforeElChildrenUpdated = options.onBeforeElChildrenUpdated || noop;
            var childrenOnly = options.childrenOnly === true;
            var fromNodesLookup = Object.create(null);
            var keyedRemovalList = [];

            function addKeyedRemoval(key) { keyedRemovalList.push(key) }

            function walkDiscardedChildNodes(node, skipKeyedNodes) {
                if (node.nodeType === ELEMENT_NODE) {
                    var curChild = node.firstChild;
                    while (curChild) {
                        var key = undefined;
                        if (skipKeyedNodes && (key = getNodeKey(curChild))) { addKeyedRemoval(key) } else { onNodeDiscarded(curChild); if (curChild.firstChild) { walkDiscardedChildNodes(curChild, skipKeyedNodes) } }
                        curChild = curChild.nextSibling
                    }
                }
            }

            function removeNode(node, parentNode, skipKeyedNodes) {
                if (onBeforeNodeDiscarded(node) === false) { return }
                if (parentNode) { parentNode.removeChild(node) }
                onNodeDiscarded(node);
                walkDiscardedChildNodes(node, skipKeyedNodes)
            }

            function indexTree(node) {
                if (node.nodeType === ELEMENT_NODE || node.nodeType === DOCUMENT_FRAGMENT_NODE$1) {
                    var curChild = node.firstChild;
                    while (curChild) {
                        var key = getNodeKey(curChild);
                        if (key) { fromNodesLookup[key] = curChild }
                        indexTree(curChild);
                        curChild = curChild.nextSibling
                    }
                }
            }
            indexTree(fromNode);

            function handleNodeAdded(el) {
                onNodeAdded(el);
                var curChild = el.firstChild;
                while (curChild) {
                    var nextSibling = curChild.nextSibling;
                    var key = getNodeKey(curChild);
                    if (key) {
                        var unmatchedFromEl = fromNodesLookup[key];
                        if (unmatchedFromEl && compareNodeNames(curChild, unmatchedFromEl)) {
                            curChild.parentNode.replaceChild(unmatchedFromEl, curChild);
                            morphEl(unmatchedFromEl, curChild)
                        } else { handleNodeAdded(curChild) }
                    } else { handleNodeAdded(curChild) }
                    curChild = nextSibling
                }
            }

            function cleanupFromEl(fromEl, curFromNodeChild, curFromNodeKey) {
                while (curFromNodeChild) {
                    var fromNextSibling = curFromNodeChild.nextSibling;
                    if (curFromNodeKey = getNodeKey(curFromNodeChild)) { addKeyedRemoval(curFromNodeKey) } else { removeNode(curFromNodeChild, fromEl, true) }
                    curFromNodeChild = fromNextSibling
                }
            }

            function morphEl(fromEl, toEl, childrenOnly) {
                var toElKey = getNodeKey(toEl);
                if (toElKey) { delete fromNodesLookup[toElKey] }
                if (!childrenOnly) {
                    if (onBeforeElUpdated(fromEl, toEl) === false) { return }
                    morphAttrs(fromEl, toEl);
                    onElUpdated(fromEl);
                    if (onBeforeElChildrenUpdated(fromEl, toEl) === false) { return }
                }
                if (fromEl.nodeName !== "TEXTAREA") { morphChildren(fromEl, toEl) } else { specialElHandlers.TEXTAREA(fromEl, toEl) }
            }

            function morphChildren(fromEl, toEl) {
                var curToNodeChild = toEl.firstChild;
                var curFromNodeChild = fromEl.firstChild;
                var curToNodeKey;
                var curFromNodeKey;
                var fromNextSibling;
                var toNextSibling;
                var matchingFromEl;
                outer: while (curToNodeChild) {
                    toNextSibling = curToNodeChild.nextSibling;
                    curToNodeKey = getNodeKey(curToNodeChild);
                    while (curFromNodeChild) {
                        fromNextSibling = curFromNodeChild.nextSibling;
                        if (curToNodeChild.isSameNode && curToNodeChild.isSameNode(curFromNodeChild)) {
                            curToNodeChild = toNextSibling;
                            curFromNodeChild = fromNextSibling;
                            continue outer
                        }
                        curFromNodeKey = getNodeKey(curFromNodeChild);
                        var curFromNodeType = curFromNodeChild.nodeType;
                        var isCompatible = undefined;
                        if (curFromNodeType === curToNodeChild.nodeType) {
                            if (curFromNodeType === ELEMENT_NODE) {
                                if (curToNodeKey) {
                                    if (curToNodeKey !== curFromNodeKey) {
                                        if (matchingFromEl = fromNodesLookup[curToNodeKey]) {
                                            if (fromNextSibling === matchingFromEl) { isCompatible = false } else {
                                                fromEl.insertBefore(matchingFromEl, curFromNodeChild);
                                                if (curFromNodeKey) { addKeyedRemoval(curFromNodeKey) } else { removeNode(curFromNodeChild, fromEl, true) }
                                                curFromNodeChild = matchingFromEl
                                            }
                                        } else { isCompatible = false }
                                    }
                                } else if (curFromNodeKey) { isCompatible = false }
                                isCompatible = isCompatible !== false && compareNodeNames(curFromNodeChild, curToNodeChild);
                                if (isCompatible) { morphEl(curFromNodeChild, curToNodeChild) }
                            } else if (curFromNodeType === TEXT_NODE || curFromNodeType == COMMENT_NODE) { isCompatible = true; if (curFromNodeChild.nodeValue !== curToNodeChild.nodeValue) { curFromNodeChild.nodeValue = curToNodeChild.nodeValue } }
                        }
                        if (isCompatible) {
                            curToNodeChild = toNextSibling;
                            curFromNodeChild = fromNextSibling;
                            continue outer
                        }
                        if (curFromNodeKey) { addKeyedRemoval(curFromNodeKey) } else { removeNode(curFromNodeChild, fromEl, true) }
                        curFromNodeChild = fromNextSibling
                    }
                    if (curToNodeKey && (matchingFromEl = fromNodesLookup[curToNodeKey]) && compareNodeNames(matchingFromEl, curToNodeChild)) {
                        fromEl.appendChild(matchingFromEl);
                        morphEl(matchingFromEl, curToNodeChild)
                    } else {
                        var onBeforeNodeAddedResult = onBeforeNodeAdded(curToNodeChild);
                        if (onBeforeNodeAddedResult !== false) {
                            if (onBeforeNodeAddedResult) { curToNodeChild = onBeforeNodeAddedResult }
                            if (curToNodeChild.actualize) { curToNodeChild = curToNodeChild.actualize(fromEl.ownerDocument || doc) }
                            fromEl.appendChild(curToNodeChild);
                            handleNodeAdded(curToNodeChild)
                        }
                    }
                    curToNodeChild = toNextSibling;
                    curFromNodeChild = fromNextSibling
                }
                cleanupFromEl(fromEl, curFromNodeChild, curFromNodeKey);
                var specialElHandler = specialElHandlers[fromEl.nodeName];
                if (specialElHandler) { specialElHandler(fromEl, toEl) }
            }
            var morphedNode = fromNode;
            var morphedNodeType = morphedNode.nodeType;
            var toNodeType = toNode.nodeType;
            if (!childrenOnly) {
                if (morphedNodeType === ELEMENT_NODE) {
                    if (toNodeType === ELEMENT_NODE) {
                        if (!compareNodeNames(fromNode, toNode)) {
                            onNodeDiscarded(fromNode);
                            morphedNode = moveChildren(fromNode, createElementNS(toNode.nodeName, toNode.namespaceURI))
                        }
                    } else { morphedNode = toNode }
                } else if (morphedNodeType === TEXT_NODE || morphedNodeType === COMMENT_NODE) { if (toNodeType === morphedNodeType) { if (morphedNode.nodeValue !== toNode.nodeValue) { morphedNode.nodeValue = toNode.nodeValue } return morphedNode } else { morphedNode = toNode } }
            }
            if (morphedNode === toNode) { onNodeDiscarded(fromNode) } else {
                if (toNode.isSameNode && toNode.isSameNode(morphedNode)) { return }
                morphEl(morphedNode, toNode, childrenOnly);
                if (keyedRemovalList) { for (var i = 0, len = keyedRemovalList.length; i < len; i++) { var elToRemove = fromNodesLookup[keyedRemovalList[i]]; if (elToRemove) { removeNode(elToRemove, elToRemove.parentNode, false) } } }
            }
            if (!childrenOnly && morphedNode !== fromNode && fromNode.parentNode) {
                if (morphedNode.actualize) { morphedNode = morphedNode.actualize(fromNode.ownerDocument || doc) }
                fromNode.parentNode.replaceChild(morphedNode, fromNode)
            }
            return morphedNode
        }
    }
    var morphdom = morphdomFactory(morphAttrs);
    return morphdom
});
/*
MORPHDOM END
*/

const MAGIC_OP_ADD = 0;
const MAGIC_OP_DEL = 1;
const MAGIC_OP_RPL = 2;
const MAGIC_OP_SWP = 3;

const magic = {
    effects: new Map(),
    currentEffect: null,
    data: reactive((() => {
        const attr = document.getElementsByTagName("html")[0].attributes;
        const ssrdata = attr.getNamedItem("data-ss");
        attr.removeNamedItem("data-ss");
        return JSON.parse(atob(ssrdata.value));
    })()),
    socket: null
}

function connect() {
    const ws_params = new URLSearchParams(location.search);
    ws_params.append("ws", document
        .getElementsByTagName("html")[0]
        .attributes.getNamedItem("data-connid").value);
    magic.socket = new WebSocket("ws://" + location.host + location.pathname + "?" + ws_params);
    magic.socket.onopen = function() {
        // // subscribe to some channels
        // magic.socket.send(JSON.stringify({
        //     //.... some message the I must send when I connect ....
        // }));
    };

    magic.socket.onmessage = function(e) {
        const msgs = JSON.parse(event.data);
        for (let i = 0; i < msgs.length; i++) {
            const msg = msgs[0];
            switch (typeof msg[0]) {
                case "number":
                    handlePatch(msg[0], msg[1], msg[2]);
                    break;
                case "string":
                    handleEvent(msg[0], msg[1]);
                    breaak;
            }
        }
    };

    magic.socket.onclose = function(e) {
        console.log('Socket is closed. Reconnect will be attempted in 1 second.', e.reason);
        setTimeout(function() {
            connect();
        }, 1000);
    };

    magic.socket.onerror = function(err) {
        console.error('Socket encountered error: ', err, 'Closing socket');
        magic.socket.close();
    };
}
connect()

function handlePatch(op, path, data) {
    switch (op) {
        case MAGIC_OP_ADD:
            set(magic.data, path, data);
            break;
        case MAGIC_OP_DEL:
            set(magic.data, path, undefined)
            break;
        case MAGIC_OP_RPL:
            set(magic.data, path, data);
            break;
        case MAGIC_OP_SWP:
            const v1 = get(magic.data, path)
            const v2 = get(magic.data, data)
            set(magic.data, path, v1)
            set(magic.data, data, v2)
            break;
    }
}

function handleEvent(name, data) {
    console.log(name, data);
}

class MagicTemplate {
    static parseTemplate(template) {
        let i = 0
        let tokens = [
            [0, ""]
        ]
        while (i < template.length) {
            if (template[i] === "{" && template[i + 1] === "§") {
                i += 2
                tokens.push([1, ""])
                while (template[i] !== "§" && template[i + 1] !== "}") {
                    tokens[tokens.length - 1][1] += template[i]
                    i += 1
                }
                i += 2
                tokens.push([0, ""])
                continue
            }
            tokens[tokens.length - 1][1] += template[i]
            i += 1
        }
        return tokens
    }

    static buildAst(tokens) {
        let ast = []
        while (tokens.length > 0) {
            if (tokens[0][0] === 1 && tokens[0][1] === "end") {
                tokens.shift()
                return ast
            } else if (tokens[0][0] === 1) {
                const parts = tokens[0][1].split(" ")
                if (parts.length > 1) {
                    tokens.shift()
                    ast.push([1, parts[0], parts[1], MagicTemplate.buildAst(tokens)])
                    continue
                }
            }
            ast.push(tokens.shift())
        }
        return ast
    }

    static exec(ast, scope, applicator) {
        let res = ""
        let i = 0
        while (i < ast.length) {
            if (ast[i][0] === 0) {
                res += ast[i][1];
            } else if (ast[i][0] === 1) {
                res += MagicTemplate.execLogic(ast[i], scope, applicator);
            }
            i++
        }
        return res
    }

    static execLogic(token, scope, applicator) {
        let res = ""
        if (token[1] === "range") {
            res += MagicTemplate.execLogicRange(token, scope, applicator)
        } else if (token[1] === "if") {
            res += MagicTemplate.execLogicShow(token, scope, applicator)
        } else {
            res += applicator(MagicTemplate.buildPath(scope, token[1]))
        }
        return res
    }

    static execLogicRange(token, scope, applicator) {
        let res = ""
        let path = MagicTemplate.buildPath(scope, token[2])
        let data = applicator(path)
        let i = 0;
        let n = Number(data)
        let amount = isNaN(n) ? data.length : n
        while (i < amount) {
            res += MagicTemplate.exec(token[3], `${path}[${i}]`, applicator)
            i++
        }
        return res
    }

    static execLogicShow(token, scope, applicator) {
        let res = ""
        let data = applicator(token[2])
        if (data !== undefined && JSON.parse(data)) {
            res += MagicTemplate.exec(token[3], scope, applicator)
        }
        return res
    }

    static buildPath(scope, path) {
        if (path === ".") {
            return scope
        } else if (path[0] === ".") {
            return scope + path
        } else {
            return path
        }
    }

    static compile(template) {
        return MagicTemplate.buildAst(MagicTemplate.parseTemplate(template))
    }
}

function reactive(object) {
    if (object === null || typeof object !== 'object') {
        return object;
    }
    for (const property in object) {
        object[property] = reactive(object[property])
    }
    return new Proxy(object, {
        get(target, property) {
            if (magic.currentEffect === null) {
                return target[property];
            }
            if (!magic.effects.has(target)) {
                magic.effects.set(target, {});
            }
            const targetEffects = magic.effects.get(target);
            if (!targetEffects[property]) {
                targetEffects[property] = [];
            }
            targetEffects[property].push(magic.currentEffect)
            return target[property];
        },
        set(target, property, value) {
            target[property] = reactive(value);
            if (magic.effects.has(target)) {
                const targetEffects = magic.effects.get(target)[property]
                if (targetEffects) {
                    let targetEffectsLength = targetEffects.length
                    let i = 0
                    while (i < targetEffectsLength) {
                        targetEffects[i]()
                        i++
                    }
                }
            }
            return true;
        },
    });
}

function get(value, path) {
    let currentEffectTemp = magic.currentEffect
    magic.currentEffect = null
    const parts = String(path).split(".");
    let acc = value;
    let target = value;
    let v = undefined;
    for (let i = 0; i < parts.length; i++) {
        v = parts[i];
        if (v.length > 0 && v[0] == "[") {
            v = Number(v.slice(1, -1));
        }
        target = acc
        acc = acc[v] !== undefined && acc[v] !== null ? acc[v] : undefined;
    }
    magic.currentEffect = currentEffectTemp
    return target[v];
}

function set(obj, keys, val) {
    let currentEffectTemp = magic.currentEffect
    magic.currentEffect = null
    keys.split && (keys = keys.split("."));
    let i = 0,
        l = keys.length,
        t = obj,
        x,
        k;
    while (i < l) {
        k = keys[i++];
        if (k === "__proto__" || k === "constructor" || k === "prototype") break;
        t = t[k] =
            i === l ?
            val :
            typeof(x = t[k]) === typeof keys ?
            x :
            keys[i] * 0 !== 0 || !!~("" + keys[i]).indexOf(".") ? {} : [];
    }
    magic.currentEffect = currentEffectTemp
}

function hydrate(element = null) {
    if (element === null) {
        return hydrate(document.children[0])
    }
    hydrater(element)
    hydrateChildren(element)
}

function hydrateChildren(element) {
    const children = element.children
    let childrenlength = children.length
    while (childrenlength--) {
        hydrate(children[childrenlength])
    }
}

function hydrater(element) {
    const attributes = element.attributes
    let attributeslen = attributes.length
    while (attributeslen--) {
        const attr = attributes[attributeslen]
        if (attr.name === "magic-value") {
            templateEffect(attr.value, (v) => {
                element.innerHTML = v
                hydrateChildren(element)
            })
            continue
        }
        if (attr.name === "magic-click") {
            element.addEventListener("click", magicClickHandler)
            continue
        }
        if (attr.name.startsWith("magic-")) {
            const template = attr.value
            const realName = attr.name.slice(6)
            templateEffect(template, (v) => element.setAttribute(realName, v))
            continue
        }
    }
}

function magicClickHandler(e) {
    const name = e.currentTarget.attributes["magic-click"].value
    const serverEvent = createServerEvent(name, {
        type: e.type,
        screenX: e.screenX,
        screenY: e.screenY,
        ctrlKey: e.ctrlKey,
        metaKey: e.metaKey,
        altKey: e.altKey,
    })
    console.log(serverEvent)
}

function createServerEvent(name, data) {
    return [name, data]
}

function effect(callback) {
    magic.currentEffect = callback;
    callback();
    magic.currentEffect = null;
}

function templateEffect(template, setter, scope = "", ) {
    if (!template.includes("{{")) {
        effect(() =>
            setter(get(magic.data, MagicTemplate.buildPath(scope, template)))
        )
        return
    }
    const tmpl = MagicTemplate.compile(template)
    effect(() =>
        setter(MagicTemplate.exec(tmpl, scope, (path) =>
            get(magic.data, MagicTemplate.buildPath(scope, path))
        ))
    )
}

document.addEventListener('DOMContentLoaded', () => hydrate(), false);