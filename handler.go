package magic

import (
	"net/http"

	"github.com/gobwas/ws"
	"github.com/klauspost/compress/gzhttp"
)

type ComponentHTTPHandler ComponentFn[Empty]

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
	renderable := f(s, Empty{})
	renderable.HTML(w)
}

func CompressedComponentHTTPHandler(fn ComponentFn[Empty]) http.Handler {
	return gzhttp.GzipHandler(ComponentHTTPHandler(fn))
}

func Compressor(h http.Handler) http.HandlerFunc {
	return gzhttp.GzipHandler(h)
}
