package magic

import (
	"encoding/json"
	"io"
	"sync"
)

// ep buffer is responsible for writing data to the client
// it holds all the data that needs to be written to it
type epBuffer struct {
	events  Events
	patches Patches
	w       io.Writer
	l       sync.Mutex
}

func (epb *epBuffer) PushEvent(name string, value any) {
	epb.l.Lock()
	epb.events.PushEvent(name, value)
	epb.l.Unlock()
	epb.Flush()
}

func (epb *epBuffer) PushPatch(op Operation, path string, value any) {
	epb.l.Lock()
	epb.patches.PushPatch(op, path, value)
	epb.l.Unlock()
	epb.Flush()
}

func (epb *epBuffer) setWriter(w io.Writer) {
	epb.l.Lock()
	epb.w = w
	epb.l.Unlock()
}

func (epb *epBuffer) Flush() {
	epb.l.Lock()
	if epb.w == nil {
		epb.l.Unlock()
		return
	}
	raw := []json.RawMessage{}
	for i := range epb.patches {
		b, _ := epb.patches[i].MarshalJSON()
		raw = append(raw, b)
	}
	for i := range epb.events {
		b, _ := epb.events[i].MarshalJSON()
		raw = append(raw, b)
	}
	epb.events = Events{}
	epb.patches = Patches{}
	b, _ := json.Marshal(raw)
	epb.w.Write(b)
	epb.l.Unlock()
}
