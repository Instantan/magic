package magic

import (
	"encoding/json"
	"log"
	"sync"
)

/*
[
	[templateID, TEMPLATE]
	[socketID, templateID, [POSITIONS], [DATA]]
]
*/

type patch struct {
	socketid  json.RawMessage
	templates []*Template
	data      map[string]any
}

type patches struct {
	p              []*patch
	l              sync.Mutex
	startedFlusher bool
	onSend         func(ps []*patch)
}

func NewPatches(onSend func(ps []*patch)) *patches {
	return &patches{
		p:      []*patch{},
		l:      sync.Mutex{},
		onSend: onSend,
	}
}

var patchPool = sync.Pool{
	New: func() any {
		return new(patch)
	},
}

func getPatch() *patch {
	return patchPool.Get().(*patch)
}

func (p *patch) free() {
	p.data = map[string]any{}
	p.socketid = []byte{}
	p.templates = []*Template{}
	patchPool.Put(p)
}

func (ps *patches) append(p ...*patch) {
	ps.l.Lock()
	ps.p = append(ps.p, p...)
	if !ps.startedFlusher {
		ps.startedFlusher = true
		go ps.runSend()
	}
	ps.l.Unlock()
}

func (ps *patches) runSend() {
	ps.l.Lock()
	cp := make([]*patch, len(ps.p))
	copy(cp, ps.p)
	ps.p = []*patch{}
	ps.l.Unlock()
	ps.onSend(cp)
}

func (s *socket) patchesToJson(ps []*patch) []byte {
	templatesToSend := []json.RawMessage{}
	dataToSend := []json.RawMessage{}
	for i := range ps {
		templateID := json.RawMessage{}
		for _, template := range ps[i].templates {
			if !s.templateIsKnown(template) {
				m := make([]json.RawMessage, 2)
				m[0], _ = json.Marshal(template.ID())
				templateID = m[0]
				m[1], _ = json.Marshal(template.String())
				t, _ := json.Marshal(m)
				templatesToSend = append(m, t)
				s.markTemplateAsKnown(template)
			}
		}
		d := make([]json.RawMessage, 3)
		d[0] = ps[i].socketid
		d[1] = templateID
		d[2], _ = json.Marshal(ps[i].data)
		data, _ := json.Marshal(d)
		dataToSend = append(dataToSend, data)
		ps[i].free()
	}
	templatesToSend = append(templatesToSend, dataToSend...)
	data, err := json.Marshal(templatesToSend)
	if err != nil {
		log.Printf("Failed sending patch: %v", err)
	}
	return data
}
