package magic

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"

	"github.com/Instantan/magic/internal"
	"github.com/puzpuzpuz/xsync"
)

// PageRenderer is a simple function that returns a structure
// the returned structure is used to render the template
type PageRenderer func(context PageContext) any

type Page struct {
	template    *Template
	renderer    PageRenderer
	connections *xsync.MapOf[string, *pageContext]
}

func CreatePage(template *Template, renderer PageRenderer) *Page {
	return &Page{
		template: template,
		renderer: renderer,

		// this holds all current and planned connections to this page
		// somehow the page context needs to get cleaned up after some time
		connections: xsync.NewMapOf[*pageContext](),
	}
}

func (p *Page) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sse, websocket := p.requestingConnection(r)
	if sse {
		p.serveHTTPSSE(w, r)
		return
	}
	if websocket {
		p.serveHTTPWebsockets(w, r)
		return
	}
	p.serveHTTPInitial(w, r)
}

func (p *Page) serveHTTPInitial(w http.ResponseWriter, r *http.Request) {
	ctx := newPageContext(p, r)
	data := p.renderer(ctx)
	isLive := false
	internal.Traverse(data, func(value any, path string) {
		if reactive, ok := value.(Reactive); ok {
			isLive = true
			reactive.Subscribe(&PatchRedirecter{
				Path:      path,
				Patchable: ctx,
			})
		}
	})

	if isLive {
		connID := randomConnectionString()
		p.connections.Store(connID, ctx)
		p.template.executeLiveSSE(w, connID, data)
		return
	}

	p.template.ExecuteStatic(w, data)
}

func (p *Page) requestingConnection(r *http.Request) (sse, websocket bool) {
	sse = r.Header.Get("Accept") == "text/event-stream"
	if sse {
		return sse, websocket
	}
	websocket = r.Header.Get("Upgrade") == "websocket"
	return sse, websocket
}

func (p *Page) serveHTTPWebsockets(w http.ResponseWriter, r *http.Request) {

}

func (p *Page) serveHTTPSSE(w http.ResponseWriter, r *http.Request) {
	connID := r.URL.Query().Get("sse")
	ctx, ok := p.connections.Load(connID)
	if !ok {
		http.Error(w, "Requested connection not available", http.StatusNotFound)
		return
	}
	s := establishSSEConnection(w)
	ctx.epb.setWriter(s)
	log.Println("Started SSE")
	<-r.Context().Done()
	ctx.epb.setWriter(nil)
	log.Println("Ended SSE")
}

func randomConnectionString() string {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return base64.RawStdEncoding.EncodeToString(randomBytes)
}
