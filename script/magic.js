import * as nanomorph from 'nanomorph'

window.magic = {
    socket: null
}

function connect() {
    const ws_params = new URLSearchParams(location.search);
    ws_params.append("ws", "0");
    window.magic.socket = new WebSocket("ws://" + location.host + location.pathname + "?" + ws_params);
    window.magic.socket.onopen = function() {
        // // subscribe to some channels
        // magic.socket.send(JSON.stringify({
        //     //.... some message the I must send when I connect ....
        // }));
    };

    window.magic.socket.onmessage = function(e) {
        console.log(e)
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