package magic

import (
	"encoding/json"
	"sync"
)

/*
[
	[templateID, TEMPLATE]
	[socketID, templateID, [POSITIONS], [DATA]]
]
*/

type patch struct {
	socketid json.RawMessage
	data     map[string]any
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
	ps.startedFlusher = false
	ps.onSend(cp)
	ps.l.Unlock()
}
