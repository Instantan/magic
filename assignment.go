package magic

import (
	"encoding/json"
	"sync"
	"time"
)

/*
a patch has the following format when it gets sent to the client
[

	[templateID, TEMPLATE]
	[socketID, templateID, DATA]

]
*/
type assignment struct {
	socketid json.RawMessage
	data     map[string]any
}

type assignments struct {
	p              []*assignment
	l              sync.Mutex
	startedFlusher bool
	onSend         func(ps []*assignment)
}

func NewPatches(onSend func(ps []*assignment)) *assignments {
	return &assignments{
		p:      []*assignment{},
		l:      sync.Mutex{},
		onSend: onSend,
	}
}

var patchPool = sync.Pool{
	New: func() any {
		return new(assignment)
	},
}

func getAssignment() *assignment {
	return patchPool.Get().(*assignment)
}

func (p *assignment) free() {
	p.data = map[string]any{}
	p.socketid = []byte{}
	patchPool.Put(p)
}

func (ps *assignments) append(p ...*assignment) {
	ps.l.Lock()
	ps.p = append(ps.p, p...)
	if !ps.startedFlusher {
		ps.startedFlusher = true
		go ps.runSend()
	}
	ps.l.Unlock()
}

func (ps *assignments) runSend() {
	time.Sleep(time.Millisecond)
	ps.l.Lock()
	cp := make([]*assignment, len(ps.p))
	copy(cp, ps.p)
	ps.p = []*assignment{}
	ps.startedFlusher = false
	ps.onSend(cp)
	ps.l.Unlock()
}
