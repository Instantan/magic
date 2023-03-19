package magic

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/puzpuzpuz/xsync"
)

// PageRenderer is a simple function that returns a structure
// the returned structure is used to render the template
type PageRenderer func(context PageContext) any

type PageContext struct {
	Page    *Page
	Request *http.Request

	sendingPayload  *sync.Mutex
	startedDebounce *atomic.Bool
	payload         *[]json.RawMessage

	connection net.Conn
}

type Page struct {
	template    *Template
	renderer    PageRenderer
	debouncer   *time.Duration
	connections *xsync.MapOf[string, *PageContext]
}

func CreatePage(template *Template, renderer PageRenderer) *Page {
	return &Page{
		template: template,
		renderer: renderer,

		// this holds all current and planned connections to this page
		connections: xsync.NewMapOf[*PageContext](),
	}
}

func (p *Page) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	println("Serve")
	// if the page uses signals it is classified as a live page
	// if it registers a event handler its classified as a websocket page
	// if not it its classifier as a sse page
	fmt.Print(r.URL)
	p.serveHTTPInitial(w, r)
}

func (p *Page) SetDebounce(duration time.Duration) {
	p.debouncer = &duration
}

func (p *Page) serveHTTPInitial(w http.ResponseWriter, r *http.Request) {
	ctx := PageContext{
		sendingPayload:  &sync.Mutex{},
		startedDebounce: &atomic.Bool{},
		payload:         &[]json.RawMessage{},
		Page:            p,
		Request:         r,
	}
	data := p.renderer(ctx)
	isLive := false
	Traverse(data, func(value any, path string) {
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
		p.connections.Store(connID, &ctx)
		data, err := json.Marshal(data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(injectDataIntoHTML(p.template.data, func() []byte {
			dataSRR := "data-ssr=\"" + base64.StdEncoding.EncodeToString(data) + "\""
			dataConnID := "data-connid=\"" + connID + "\""
			return []byte(" " + dataSRR + " " + dataConnID)
		}, injectSSEScript))
	} else {
		p.template.ExecuteStatic(w, data)
	}
}

func (p *Page) serveHTTPWebsockets(w http.ResponseWriter, r *http.Request) {

}

func (p *Page) serveHTTPSSE(w http.ResponseWriter, r *http.Request) {

}

func (p PageContext) Patch(patch Patch) {
	b, err := patch.MarshalJSON()
	if err != nil {
		log.Println(err)
		return
	}
	p.send([]json.RawMessage{b})
}

func (p PageContext) send(data []json.RawMessage) {
	if p.Page.debouncer == nil {
		*p.payload = data
		p.sendPayload()
		return
	} else {
		p.sendingPayload.Lock()
		*p.payload = append(*p.payload, data...)
		p.sendingPayload.Unlock()
		if !p.startedDebounce.Load() {
			go func() {
				time.Sleep(*p.Page.debouncer)
				p.sendingPayload.Lock()
				p.sendPayload()
				p.sendingPayload.Unlock()
				p.startedDebounce.Store(false)
			}()
		}
	}
}

func (p *PageContext) sendPayload() error {
	b, err := json.Marshal(p.payload)
	if err != nil {
		return err
	}
	print(string(b))
	*p.payload = []json.RawMessage{}
	return nil
}

func randomConnectionString() string {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return base64.RawStdEncoding.EncodeToString(randomBytes)
}
