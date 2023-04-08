package magic

import _ "embed"

//go:embed script/magic-reactivity.js
var ReactivityScript []byte

func init() {
	ReactivityScript = append(append([]byte("<script>"), ReactivityScript...), []byte("</script>")...)
}

func injectReactivityScript() []byte {
	return ReactivityScript
}
