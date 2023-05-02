package magic

import (
	"context"
	"net/http"

	"github.com/gobwas/ws"
)

type HandlerFunc ComponentFn

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s := &socket{
		conn:           nil,
		ctx:            r.Context(),
		knownTemplates: NewSet[int](),
	}
	if r.Header.Get("Upgrade") == "websocket" {
		s.ctx = context.Background()
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			return
		}
		go s.establishConnection(ComponentFn(f), conn)
		return
	}
	renderable := f(s)
	renderable.HTML(w)
}
