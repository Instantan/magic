const MAGIC_OP_ADD = 0
const MAGIC_OP_DEL = 1
const MAGIC_OP_RPL = 2
const MAGIC_OP_SWP = 3

const _internal_magic = {
    _connid: document.getElementsByTagName("html")[0].attributes.getNamedItem("data-connid").value,
}

const magic = {
    source: connectToSSE(),
    data: readSSRData()
}

function connectToSSE() {
    const sse_params = new URLSearchParams(location.search)
    sse_params.append("sse", _internal_magic._connid)
    const source = new EventSource(location.pathname + "?" + sse_params);
    source.onopen = function(event) {
        console.log("Connected:", event)
    }
    source.onmessage = function(event) {
        const msgs = JSON.parse(event.data)
        for (let i = 0; i < msgs.length; i++) {
            const msg = msgs[0]
            switch (typeof(msg[0])) {
                case 'number':
                    handlePatch(msg[0], msg[1], msg[2])
                    break
                case 'string':
                    handleEvent(msg[0], msg[1])
                    breaak
            }
        }
    }
    source.onerror = function(event) {
        console.error(event)
    }
    return source
}

function readSSRData() {
    const attr = document.getElementsByTagName("html")[0].attributes
    const ssrdata = attr.getNamedItem("data-ssr")
    attr.removeNamedItem("data-ssr")
    return JSON.parse(atob(ssrdata.value))
}

function handlePatch(op, path, data) {
    console.log(op, path, data)
}

function handleEvent(name, data) {
    console.log(name, data)
}