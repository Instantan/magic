package magic

import (
	"net/http"

	"github.com/gobwas/ws"
	"github.com/klauspost/compress/gzhttp"
)

type Options struct {
	// Injects the magic script if its empty, if not it adds a defered script src with the given url
	MagicScriptURL string
	// OnlyStatic disables the websocket (live) connection when true
	OnlyStatic bool
	// Compressed enables gzip compression for the handler
	Compressed bool
}

type Option func(opts *Options)

func WithOptions(options Options) Option {
	return func(opts *Options) {
		*opts = options
	}
}

func WithOnlyStatic() Option {
	return func(opts *Options) {
		opts.OnlyStatic = true
	}
}

func WithCompressed() Option {
	return func(opts *Options) {
		opts.Compressed = true
	}
}

func ComponentHTTPHandler(fn ComponentFn[Empty], options ...Option) http.Handler {
	opts := &Options{}
	for _, optFn := range options {
		optFn(opts)
	}
	config := &htmlRenderConfig{
		static: opts.OnlyStatic,
	}
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s := NewSocket(r)
		if r.Header.Get("Upgrade") == "websocket" {
			r.Header.Set(gzhttp.HeaderNoCompression, "-")
			conn, _, _, err := ws.UpgradeHTTP(r, w)
			if err != nil {
				return
			}
			submitTask(func() {
				s.establishConnection(ComponentFn[Empty](fn), conn)
			})
			return
		}
		av := fn(s, Empty{})
		s.deferredAssigns.Wait()
		av.html(w, config)
	})
	if opts.Compressed {
		return gzhttp.GzipHandler(handler)
	} else {
		return handler
	}
}
