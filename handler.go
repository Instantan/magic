package magic

import (
	"net/http"

	"github.com/gobwas/ws"
	"github.com/klauspost/compress/gzhttp"
)

type Options struct {
	// OnlyStatic disables the websocket (live) connection when true
	OnlyStatic bool
	// Compressed enables gzip compression for the handler
	Compressed bool
	// Injects the magic script if its empty, if not it adds a defered script src with the given url
	MagicScriptURL string
}

type Option func(opts *Options)

func WithOptions(options Options) Option {
	return func(opts *Options) {
		*opts = options
	}
}

func WithOnlyStatic(onlyStatic bool) Option {
	return func(opts *Options) {
		opts.OnlyStatic = onlyStatic
	}
}

func WithCompressed(compressed bool) Option {
	return func(opts *Options) {
		opts.Compressed = compressed
	}
}

func WithMagicScriptURL(url string) Option {
	return func(opts *Options) {
		opts.MagicScriptURL = url
	}
}

func ComponentHTTPHandler(fn ComponentFn[Empty], options ...Option) http.Handler {
	opts := &Options{}
	for _, optFn := range options {
		optFn(opts)
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
		av.html(w, &htmlRenderConfig{
			magicScriptInline: opts.MagicScriptURL == "",
			magicScriptUrl:    opts.MagicScriptURL,
			static:            opts.OnlyStatic,
		})
	})
	if opts.Compressed {
		return gzhttp.GzipHandler(handler)
	} else {
		return handler
	}
}
