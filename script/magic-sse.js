const MAGIC_OP_ADD = 0;
const MAGIC_OP_DEL = 1;
const MAGIC_OP_RPL = 2;
const MAGIC_OP_SWP = 3;

class MagicReactor {

    constructor() {
        this.numberOfSubscribers = 0;
        this.subscribers = {}
    }

    subscribe = (path, fn) => {
        const parts = path.split(".")
        let parent = this.subscribers
        for (let i = 0; i < parts.length; i++) {
            let v = parts[i]
            if (v.length > 0 && v[0] == "[") {
                v = Number(v.slice(1, -1));
            }
            if (parent[v] === undefined) {
                parent[v] = {
                    "$subscribers": [],
                }
            }
            parent = parent[v]
        }
        parent["$subscribers"].push(fn)
        this.numberOfSubscribers++;
    }

    unsubscribe = (path, fn) => {
        const parts = path.split(".")
        let parent = this.subscribers
        for (let i = 0; i < parts.length; i++) {
            let v = parts[i]
            if (v.length > 0 && v[0] == "[") {
                v = Number(v.slice(1, -1));
            }
            if (parent[v] === undefined) {
                return
            }
            parent = parent[v]
        }
        const index = parent["$subscribers"].indexOf(fn);
        if (index > -1) {
            parent["$subscribers"].splice(index, 1);
            this.numberOfSubscribers--;
        }
    }

    notify = (path) => {
        const parts = path.split(".")
        let parent = this.subscribers
        for (let i = 0; i < parts.length; i++) {
            let v = parts[i]
            if (v.length > 0 && v[0] == "[") {
                v = Number(v.slice(1, -1));
            }
            if (parent[v] === undefined) {
                return
            }
            parent = parent[v]
        }
        this.notifyAllOf(parent)
    }

    notifyAllOf = (obj) => {
        const keys = Object.keys(obj)
        let keyslen = keys.length
        while (keyslen--) {
            if (keys[keyslen] === "$subscribers") {
                const subs = obj[keys[keyslen]]
                let subslen = subs.length
                while (subslen--) {
                    subs[subslen]()
                }
                continue
            }
            this.notifyAllOf(obj[keys[keyslen]])
        }
    }

    debug() {
        console.log("Number of subscribers", this.numberOfSubscribers)
    }
}

class MagicValue extends HTMLElement {
    subscribed = false
    absolutePath = ""

    constructor() {
        super();
    }

    recomputePath = () => {
        this.setAttribute("tag", this.attributes.tag === undefined ? this.innerHTML : this.attributes.tag.value)
        this.absolutePath = absolutePath(this, this.attributes.tag.value);
    }

    connectedCallback() {
        this.recomputePath()
        this.rerender()
    }

    sub = () => {
        if (this.subscribed) {
            _internal_magic._reactor.unsubscribe(this.absolutePath, this.rerender);
        }
        _internal_magic._reactor.subscribe(this.absolutePath, this.rerender);
        this.subscribed = true
    }

    rerender = () => {
        if (!this.subscribed) {
            this.sub()
        }
        this.innerHTML = MagicUtil.get(magic.data, this.absolutePath, this.attributes.tag.value);
    }

    disconnectedCallback() {
        _internal_magic._reactor.unsubscribe(this.absolutePath, this.rerender);
    }
}

class MagicScope extends HTMLElement {
    constructor() {
        super();
    }
    connectedCallback() {
        registerChilds(this)
    }
    disconnectedCallback() {}
}

class MagicRange extends HTMLElement {
    template = undefined;
    subscribed = false;
    constructor() {
        super();
    }
    connectedCallback() {
        this.recomputePath()
        this.rerender()
    }
    attributeChangedCallback() {
        this.recomputePath()
        this.sub()
        this.rerender()
    }
    recomputePath = () => {
        this.absolutePath = absolutePath(this, MagicTemplate.unpackTemplate(this.attributes.of.value));
    }
    sub = () => {
        if (this.subscribed) {
            _internal_magic._reactor.unsubscribe(this.absolutePath, this.rerender);
        }
        _internal_magic._reactor.subscribe(this.absolutePath, this.rerender);
        this.subscribed = true;
    }
    rerender() {
        if (!this.template) {
            const tmpl = this.innerHTML;
            this.template = (scope) => `<m-scope m-scope="${scope}">${tmpl}</m-scope>`;
        }
        if (!this.subscribed) {
            this.sub()
        }
        let path = this.attributes.of.value
        if (!MagicTemplate.isTemplate(path)) {
            return
        }
        path = absolutePath(this, this.absolutePath)
        const data = MagicUtil.get(magic.data, path, [])
        if (!Array.isArray(data)) {
            return
        }
        let l = data.length;
        this.innerHTML = ""
        while (l--) {
            this.innerHTML = this.template(path + `.[${l}]`) + this.innerHTML
        }
    }
    disconnectedCallback() {
        _internal_magic._reactor.unsubscribe(this.absolutePath, this.rerender);
    }
    static get observedAttributes() {
        return ['of']
    }
}

class MagicShow extends HTMLElement {
    template = undefined;
    subscribed = false;
    previous = false;
    absolutePath = "";
    constructor() {
        super();
    }
    recomputePath = () => {
        this.absolutePath = absolutePath(this, MagicTemplate.unpackTemplate(this.attributes.when.value));
    }
    connectedCallback() {
        this.recomputePath()
        this.rerender()
    }
    attributeChangedCallback() {
        this.recomputePath()
        this.sub()
        this.rerender()
    }
    sub = () => {
        if (this.subscribed) {
            _internal_magic._reactor.unsubscribe(this.absolutePath, this.rerender);
        }
        _internal_magic._reactor.subscribe(this.absolutePath, this.rerender);
        this.subscribed = true;
    }
    rerender = () => {
        if (!this.template) {
            this.template = this.innerHTML;
        }
        let w = this.attributes.when.value;
        if (MagicTemplate.isTemplate(w)) {
            if (!this.subscribed) {
                this.sub()
            }
            w = MagicUtil.get(magic.data, this.absolutePath, "false")
        }
        if (w === "false" || w === "0" || w === "null" || w === "undefined" ? false : Boolean(w) && this.previous === false) {
            this.innerHTML = this.template;
            this.previous = true;
            registerChilds(this);
        } else if (this.previous === true) {
            this.innerHTML = ''
            this.previous = false;
        }
    }
    disconnectedCallback() {
        _internal_magic._reactor.unsubscribe(this.absolutePath, this.rerender);
    }
    static get observedAttributes() {
        return ['when']
    }
}

class MagicUtil {

    static get(value, path, defaultValue) {
        const parts = String(path).split(".");
        let acc = value;
        for (let i = 0; i < parts.length; i++) {
            let v = parts[i];
            if (v.length > 0 && v[0] == "[") {
                v = Number(v.slice(1, -1));
            }
            acc = acc[v] !== undefined && acc[v] !== null ? acc[v] : defaultValue;
        }
        return acc;
    }

    static set(obj, keys, val) {
        keys.split && (keys = keys.split("."));
        var i = 0,
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
    }

    static connectToSSE() {
        const sse_params = new URLSearchParams(location.search);
        sse_params.append("sse", _internal_magic._connid);
        const source = new EventSource(location.pathname + "?" + sse_params);
        source.onopen = function(event) {
            console.log("Connected:", magic.data);
        };
        source.onmessage = function(event) {
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
        source.onerror = function(event) {
            console.error(event);
        };
        return source;
    }

    static readSSData() {
        const attr = document.getElementsByTagName("html")[0].attributes;
        const ssrdata = attr.getNamedItem("data-ss");
        attr.removeNamedItem("data-ss");
        return JSON.parse(atob(ssrdata.value));
    }
}

class MagicTemplate {

    static parseTemplate(template) {
        let i = 0
        let tokens = [
            [0, ""]
        ]
        while (i < template.length) {
            if (template[i] === "{" && template[i + 1] === "{") {
                i += 1
                tokens.push([1, ""])
                while (template[i] !== "}" && template[i + 1] !== "}") {
                    i += 1
                    tokens[tokens.length - 1][1] += template[i]
                }
                i += 3
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
        } else if (token[1] === "show") {
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
        while (i < data.length) {
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

    static isTemplate(template = "") {
        return template.length > 4 && template.includes("{{") && template.includes("}}")
    }

    static unpackTemplate(template = "") {
        return template.replace("{{", "").replace("}}", "")
    }

    static paths(ast) {
        const set = new Set()
        MagicTemplate.exec(ast, "", (path) => {
            set.add(path)
            return "true"
        })
        return [...set.values()]
    }
}


const _internal_magic = {
    _connid: document
        .getElementsByTagName("html")[0]
        .attributes.getNamedItem("data-connid").value,
    _reactor: new MagicReactor()
};

const magic = {
    source: MagicUtil.connectToSSE(),
    data: MagicUtil.readSSData(),
};

function handlePatch(op, path, data) {
    // The notify is technically not correct when patching a deep object
    // for example when running 
    // RPL "data.deep" { "user": { "name": { "prename": "Paul", "surename": "Blob" } } }
    // then a effect that is listening to "data.deep" gets notified but not 
    // a effect that is listening to "data.deep.user.name.prename"
    // the listening system needs to be more granular
    switch (op) {
        case MAGIC_OP_ADD:
            MagicUtil.set(magic.data, path, data);
            break;
        case MAGIC_OP_DEL:
            MagicUtil.set(magic.data, path, undefined)
            break;
        case MAGIC_OP_RPL:
            MagicUtil.set(magic.data, path, data);
            break;
        case MAGIC_OP_SWP:
            const v1 = MagicUtil.get(magic.data, path)
            const v2 = MagicUtil.get(magic.data, data)
            MagicUtil.set(magic.data, path, v1)
            MagicUtil.set(magic.data, data, v2)
            break;
    }
    _internal_magic._reactor.notify(path);
}

function handleEvent(name, data) {
    console.log(name, data);
}

function registerNodeAndChilds(node) {
    if (!node.tagName.startsWith("M-")) {
        const scope = scopeFromNode(node)
        const dataGetter = (path) => {
            return MagicUtil.get(magic.data, path, "")
        }
        let al = node.attributes.length
        while (al--) {
            const attr = node.attributes[al]
            if (MagicTemplate.isTemplate(attr.value)) {
                const template = MagicTemplate.compile(attr.value)
                MagicTemplate.paths(template).forEach(p => {
                    _internal_magic._reactor.subscribe(MagicTemplate.buildPath(scope, p), () => {
                        node.setAttribute(attr.name, MagicTemplate.exec(template, scope, dataGetter))
                    })
                })
            }
        }
        let cl = node.children.length;
        while (cl--) {
            registerNodeAndChilds(node.children[cl])
        }
    }
}

function registerChilds(node) {
    let cl = node.children.length;
    while (cl--) {
        registerNodeAndChilds(node.children[cl])
    }
}

function getScope(node) {
    let n = node
    let scope = ""
    while (true) {
        n = n.parentElement
        if (!n) {
            return scope
        }
        const scp = scopeFromNode(n)
        if (scp === "") {
            continue
        }
        scope = scp + scope
        if (!isRelativePath(scp)) {
            return scope
        }
    }
    return scope
}

function isRelativePath(path) {
    return path[0] === "."
}

function scopeFromNode(node) {
    const scopeattr = node.attributes["m-scope"]
    if (!scopeattr) {
        return ""
    }
    return scopeattr.value
}

function absolutePath(node, path) {
    if (path[0] === ".") {
        const scope = getScope(node)
        if (path === ".") {
            return scope
        }
        return scope + path
    }
    return path
}

document.addEventListener("DOMContentLoaded", function(event) {
    registerNodeAndChilds(document.getElementsByTagName("html")[0])
    customElements.define("m-v", MagicValue);
    customElements.define("m-show", MagicShow);
    customElements.define("m-range", MagicRange);
    customElements.define("m-scope", MagicScope);
});