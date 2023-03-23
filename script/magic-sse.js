const MAGIC_OP_ADD = 0;
const MAGIC_OP_DEL = 1;
const MAGIC_OP_RPL = 2;
const MAGIC_OP_SWP = 3;

const _internal_magic = {
  _connid: document
    .getElementsByTagName("html")[0]
    .attributes.getNamedItem("data-connid").value,
  _subscribers: {},
};

const magic = {
  source: connectToSSE(),
  data: readSSRData(),
};

function connectToSSE() {
  const sse_params = new URLSearchParams(location.search);
  sse_params.append("sse", _internal_magic._connid);
  const source = new EventSource(location.pathname + "?" + sse_params);
  source.onopen = function (event) {
    console.log("Connected:", magic.data);
  };
  source.onmessage = function (event) {
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
  source.onerror = function (event) {
    console.error(event);
  };
  return source;
}

function readSSRData() {
  const attr = document.getElementsByTagName("html")[0].attributes;
  const ssrdata = attr.getNamedItem("data-ssr");
  attr.removeNamedItem("data-ssr");
  return JSON.parse(atob(ssrdata.value));
}

function handlePatch(op, path, data) {
  console.log(op, path, data);
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
  notify(path);
}

function handleEvent(name, data) {
  console.log(name, data);
}

function get(value, path, defaultValue) {
  const parts = String(path).split(".");
  let acc = value;
  for (let i = 0; i < parts.length; i++) {
    const v = parts[i];
    if (v.length > 0 && v[0] == "[") {
      v = Number(v.slice(1, -1));
    }
    acc = acc[v] !== undefined && acc[v] !== null ? acc[v] : defaultValue;
  }
  return acc;
}

function set(obj, keys, val) {
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
      i === l
        ? val
        : typeof (x = t[k]) === typeof keys
        ? x
        : keys[i] * 0 !== 0 || !!~("" + keys[i]).indexOf(".")
        ? {}
        : [];
  }
}

function subscribe(path, callback) {
  if (!_internal_magic._subscribers[path]) {
    _internal_magic._subscribers[path] = [];
  }
  _internal_magic._subscribers[path].push(callback);
}

function unsubscribe(path, callback) {
  if (!_internal_magic._subscribers[path]) {
    return;
  }
  const index = _internal_magic._subscribers[path].indexOf(callback);
  if (index > -1) {
    _internal_magic._subscribers[path].splice(index, 1); // 2nd parameter means remove one item only
  }
}

function notify(path) {
  let elms = _internal_magic._subscribers[path];
  if (elms) {
    let l = elms.length;
    while (l--) {
      elms[l]();
    }
  }
}

document.addEventListener("DOMContentLoaded", function (event) {
  customElements.define(
    "m-v",
    class extends HTMLElement {
      template = undefined;
      notify = undefined;
      constructor() {
        super();
      }
      connectedCallback() {
        if (!this.template) {
          this.template = this.innerHTML;
          _internal_magic._subscribers;
          this.notify = () => this.innerHTML = get(magic.data, this.template);
          subscribe(this.template, this.notify);
        }
        this.innerHTML = get(magic.data, this.template);
      }
      disconnectedCallback() {
        unsubscribe(this.template, this.notify);
      }
    }
  );
});
