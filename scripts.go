package magic

import _ "embed"

//go:embed script/magic-sse.js
var SSEScript []byte

//go:embed script/magic-websocket.js
var WebsocketScript []byte

func init() {
	SSEScript = append(append([]byte("<script>"), SSEScript...), []byte("</script>")...)
	WebsocketScript = append(append([]byte("<script>"), WebsocketScript...), []byte("</script>")...)
}

func injectSSEScript() []byte {
	return SSEScript
}

func injectWebsocketScript() []byte {
	return WebsocketScript
}
