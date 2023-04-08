package magic

import (
	"crypto/rand"
	"encoding/base64"
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
	if r.Header.Get("Upgrade") == "websocket" {
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
		p.template.executeLiveTemplate(w, connID, data)
		return
	}

	p.template.ExecuteStatic(w, data)
}

func (p *Page) serveHTTPWebsockets(w http.ResponseWriter, r *http.Request) {
	connID := r.URL.Query().Get("ws")
	ctx, ok := p.connections.Load(connID)
	if !ok {
		http.Error(w, "Requested connection not available", http.StatusNotFound)
		return
	}
	s := establishWSConnection(w, r)
	go s.Read(ctx)
}

func randomConnectionString() string {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return base64.RawStdEncoding.EncodeToString(randomBytes)
}
