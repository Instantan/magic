package magic

import (
	"net/http"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type HandlerFunc ComponentFn

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s := &socket{
		conn:           nil,
		knownTemplates: NewSet[int](),
	}
	if r.Header.Get("Upgrade") == "websocket" {
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		s.patches = NewPatches(s.onSendTemplatePatch)
		if err != nil {
			return
		}
		renderable := f(s)
		patches := renderable.Patch()
		s.conn = conn
		data := s.patchesToJson(patches)
		wsutil.WriteServerText(s.conn, data)
		return
	}
	renderable := f(s)
	renderable.HTML(w)
}
