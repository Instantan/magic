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