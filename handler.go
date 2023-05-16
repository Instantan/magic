package magic

import (
	"net/http"

	"github.com/gobwas/ws"
)

type HandlerFunc ComponentFn[Empty]

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s := NewSocket(r)
	if r.Header.Get("Upgrade") == "websocket" {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			return
		}
		go s.establishConnection(ComponentFn[Empty](f), conn)
		return
	}
	renderable := f(s, Empty{})
	renderable.HTML(w)
}
