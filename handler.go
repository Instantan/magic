package magic

import (
	"net/http"

	"github.com/gobwas/ws"
)

type HandlerFunc ComponentFn

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s := &socket{}
	renderable := f(s)
	renderable.HTML(w)
}

func (f HandlerFunc) websocket(w http.ResponseWriter, r *http.Request) {
	conn, _, _, err := ws.UpgradeHTTP(r, w)
}
