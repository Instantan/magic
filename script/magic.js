const magic = {
    data: new Map(),
    updaters: new Map(),
    templates: new Map(),
}

function renderTemplate(template = "") {
    let t = magic.templates.get(template);
    let i = 0;
    let tl = t.length;
    let rendered = "";
    while (i < tl) {
        if (isTag(t[i])) {
            rendered += magic.data.get(t[i].slice(2, -2))
        } else {
            rendered += t[i]
        }
        i++
    }
    return rendered
}

function compileTemplate(template = "") {
    if (magic.templates.has(template)) {
        let templateValues = []
        let t = magic.templates.get(template)
        let i = t.length
        while (i--) {
            if (isTag(t[i])) {
                templateValues.push(t[i].slice(2, -2))
            }
        }
        return templateValues
    }
    let templateLength = template.length;
    let cold = false
    let compiled = [""];
    let templateValues = []
    for (let i = 0; i < templateLength; i++) {
        let c = template[i]
        if (!cold && c == "{" && i + 1 < templateLength && template[i + 1] === "{") {
            cold = true
            compiled.push("{{")
            i++
        } else if (cold && c == "}" && i + 1 < templateLength && template[i + 1] === "}") {
            cold = false
            compiled[compiled.length - 1] += "}}"
            templateValues.push(compiled[compiled.length - 1].slice(2, -2))
            compiled.push("")
            i++
        } else {
            compiled[compiled.length - 1] += c
        }
    }
    magic.templates.set(template, compiled)
    return templateValues
}

function isTag(template = "") {
    return template.length > 4 && template.includes("{{") && template.includes("}}")
}

function registerAttributeTemplatesForCollection(collection) {
    let i = collection.length
    while (i--) {
        let item = collection[i]
        registerAttributeTemplatesForItem(item)
        registerAttributeTemplatesForCollection(item.children)
    }
}

function registerAttributeTemplatesForItem(item) {
    let i = item.attributes.length
    while (i--) {
        let attribute = item.attributes[i]
        if (isTag(attribute.value)) {
            let av = attribute.value
            const templateNames = compileTemplate(av)
            let itnl = templateNames.length
            while (itnl--) {
                registerUpdater(templateNames[itnl], () => item.setAttribute(attribute.name, renderTemplate(av)))
            }
        }
    }
}

function registerUpdater(key, updater) {
    if (!magic.updaters.has(key)) {
        magic.updaters.set(key, [updater])
    } else {
        magic.updaters.set(key, magic.updaters.get(key).push(updater))
    }
}

function runUpdaters(key) {
    let updaters = magic.updaters.get(key)
    if (updaters) {
        let i = updaters.length
        while (i--) {
            updaters[i]()
        }
    }
}

class MagicValue extends HTMLElement {
    constructor() {
        super();
        this.name = this.innerHTML;
        registerUpdater(this.innerHTML, () => this.innerHTML = magic.data.get(this.name))
        this.innerHTML = magic.data.get(this.name)
    }
}

document.addEventListener("DOMContentLoaded", function(event) {
    document.body.innerHTML = replaceTemplatesInInnerHTMLWithMagicValue(document.body.innerHTML);
    registerAttributeTemplatesForCollection(document.children)
    customElements.define("m-v", MagicValue);
});