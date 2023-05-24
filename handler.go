package magic

import (
	"net/http"

	"github.com/gobwas/ws"
	"github.com/klauspost/compress/gzhttp"
)

type ComponentHTTPHandler ComponentFn[Empty]
type StaticComponentHTTPHandler ComponentFn[Empty]

func (f ComponentHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s := NewSocket(r)
	if r.Header.Get("Upgrade") == "websocket" {
		r.Header.Set(gzhttp.HeaderNoCompression, "-")
		conn, _, _, err := ws.UpgradeHTTP(r, w)
		if err != nil {
			return
		}
		go s.establishConnection(ComponentFn[Empty](f), conn)
		return
	}
	f(s, Empty{}).html(w, &htmlRenderConfig{
		magicScriptInline: true,
	})
}

func (f StaticComponentHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(NewSocket(r), Empty{}).html(w, &htmlRenderConfig{
		static: true,
	})
}

func CompressedComponentHTTPHandler(fn ComponentFn[Empty]) http.Handler {
	return gzhttp.GzipHandler(ComponentHTTPHandler(fn))
}

func CompressedStaticComponentHTTPHandler(fn ComponentFn[Empty]) http.Handler {
	return gzhttp.GzipHandler(StaticComponentHTTPHandler(fn))
}
