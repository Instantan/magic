package magic

import (
	"net/http"

	"github.com/gobwas/ws"
)

type HandlerFunc ComponentFn[Empty]

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s := &socket{
		conn:           nil,
		socketrefs:     map[uintptr]Socket{},
		socketrefsRefs: map[uintptr]uint{},
		knownTemplates: NewSet[int](),
	}
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
