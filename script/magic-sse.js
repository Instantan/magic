const magic = {
    _connid: document.getElementsByTagName("html")[0].attributes.getNamedItem("data-connid").value
}
const sse_params = new URLSearchParams(location.search).append("connid", magic._connid)
const source = new EventSource(location.pathname + sse_params);
source.onmessage = function(event) {
    document.getElementById("result").innerHTML += event.data + "<br>";
}