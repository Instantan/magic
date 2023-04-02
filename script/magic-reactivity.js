/*
<html lang="en">

<head>
    <title magic-value="nested.data">Mikado Basic Example</title>
</head>

<script>
    const magic = {
        effects: new Map(),
        currentEffect: null,
        data: reactive({
            "nested": {
                "data": 1,
                "bla": [2]
            }
        })
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
                    let targetEffectsLength = targetEffects.length
                    let i = 0
                    while (i < targetEffectsLength) {
                        targetEffects[i]()
                        i++
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

    function hydrate(element = null) {
        if (element === null) {
            return hydrate(document.children[0])
        }
        hydrater(element)
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
                effect(() => element.innerHTML = get(magic.data, attr.value))
                continue
            }
            if (attr.name === "magic-click") {
                element.addEventListener("click", magicClickHandler)
                continue
            }
            if (attr.name.startsWith("magic-")) {
                const template = attr.value
                const realName = attr.name.slice(6)
                effect(() => element.setAttribute(realName, get(magic.data, attr.value)))
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
</script>

<body magic-class="nested.data">

    <button magic-click="test">
        <v magic-value="nested.data"></v>
    </button>

</body>

<script>
    function effect(callback) {
        magic.currentEffect = callback;
        callback();
        magic.currentEffect = null;
    }


    // effect(() => {
    //     console.log(get(magic.data, "nested.data"))
    // })

    setInterval(() => {
        magic.data.nested.data++
    }, 1000)

    hydrate()
</script>

</html>
*/

const magic = {
    effects: new Map(),
    currentEffect: null,
    data: reactive({
        "nested": {
            "data": 1,
            "bla": [2]
        }
    })
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
                let targetEffectsLength = targetEffects.length
                let i = 0
                while (i < targetEffectsLength) {
                    targetEffects[i]()
                    i++
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

function hydrate(element = null) {
    if (element === null) {
        return hydrate(document.children[0])
    }
    hydrater(element)
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
            effect(() => element.innerHTML = get(magic.data, attr.value))
            continue
        }
        if (attr.name === "magic-click") {
            element.addEventListener("click", magicClickHandler)
            continue
        }
        if (attr.name.startsWith("magic-")) {
            const template = attr.value
            const realName = attr.name.slice(6)
            effect(() => element.setAttribute(realName, get(magic.data, attr.value)))
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