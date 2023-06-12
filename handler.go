package magic

import (
	"net/http"

	"github.com/gobwas/ws"
	"github.com/klauspost/compress/gzhttp"
)

type ComponentHTTPHandler ComponentFn[Empty]
type StaticComponentHTTPHandler ComponentFn[Empty]

type Server interface {
	ComponentHTTPHandler(fn ComponentFn[Empty]) http.Handler
	CompressedComponentHTTPHandler(fn ComponentFn[Empty]) http.Handler
	StaticComponentHTTPHandler(fn ComponentFn[Empty]) http.Handler
	CompressedStaticComponentHTTPHandler(fn ComponentFn[Empty]) http.Handler
	MagicScriptHandler(http.ResponseWriter, *http.Request)
}

type deferedMagicScriptServer struct {
	url string
}

func NewServer(magicScriptUrl string) Server {
	return deferedMagicScriptServer{url: magicScriptUrl}
}

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
	av := f(s, Empty{})
	s.deferredAssigns.Wait()
	av.html(w, &htmlRenderConfig{
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

func (dmss deferedMagicScriptServer) ComponentHTTPHandler(fn ComponentFn[Empty]) http.Handler {
	return dmss.componentHTTPHandler(fn)
}

func (dmss deferedMagicScriptServer) CompressedComponentHTTPHandler(fn ComponentFn[Empty]) http.Handler {
	return gzhttp.GzipHandler(dmss.componentHTTPHandler(fn))
}

func (dmss deferedMagicScriptServer) StaticComponentHTTPHandler(fn ComponentFn[Empty]) http.Handler {
	return StaticComponentHTTPHandler(fn)
}

func (dmss deferedMagicScriptServer) CompressedStaticComponentHTTPHandler(fn ComponentFn[Empty]) http.Handler {
	return CompressedStaticComponentHTTPHandler(fn)
}

func (dmss deferedMagicScriptServer) componentHTTPHandler(f ComponentFn[Empty]) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		av := f(s, Empty{})
		s.deferredAssigns.Wait()
		av.html(w, &htmlRenderConfig{
			magicScriptInline: false,
			magicScriptUrl:    dmss.url,
		})
	})
}

func (dmss deferedMagicScriptServer) MagicScriptHandler(w http.ResponseWriter, r *http.Request) {
	ServeMagicScript(w, r)
}
