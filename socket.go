package magic

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type Socket interface {
	Live() bool
	DispatchEvent(ev string, data any) error
	HandleEvent(EventHandler)
	Request() *http.Request

	id() uintptr
	clone() Socket
	assign(key string, value any)
	track(Socket)
	untrack(Socket)
	dispatch(ev string, data EventData)
}

type countableSocket struct {
	socket Socket
	count  int
}

type socket struct {
	refs           map[uintptr]countableSocket
	knownTemplates set[int]

	conn        net.Conn
	assignments *assignments

	sending  sync.Mutex
	tracking sync.Mutex

	request         *http.Request
	deferredAssigns sync.WaitGroup
	initialized     bool
}

func NewSocket(request *http.Request) *socket {
	return &socket{
		conn:           nil,
		refs:           map[uintptr]countableSocket{},
		request:        request,
		knownTemplates: NewSet[int](),
	}
}

func (s *socket) Live() bool {
	return s.assignments != nil
}

func (s *socket) HandleEvent(_ EventHandler) {
	// we need to dispatch the event to the right ref
}

func (s *socket) Request() *http.Request {
	return s.request
}

func (s *socket) handleEvent(ev Event) {
	if ev.Target == 0 {
		s.dispatch(ev.Kind, EventData(ev.Payload))
		return
	}
	sref, ok := s.refs[ev.Target]
	if !ok {
		return
	}
	sref.socket.dispatch(ev.Kind, EventData(ev.Payload))
}

func (s *socket) DispatchEvent(ev string, data any) error {
	if !s.Live() {
		return nil
	}
	s.dispatchEvent(ev, data, s.id())
	return nil
}

func (s *socket) dispatchEvent(ev string, data any, target uintptr) error {
	m, err := json.Marshal(data)
	if err != nil {
		return err
	}
	d, err := json.Marshal([]Event{
		{
			Kind:    ev,
			Target:  target,
			Payload: m,
		},
	})
	if err != nil {
		return err
	}
	s.send(d)
	return nil
}

func (s *socket) id() uintptr {
	return 0
}

func (s *socket) clone() Socket {
	return &ref{
		root:  s,
		state: map[string]any{},
	}
}

func (s *socket) assign(key string, value any) {
}

func (s *socket) establishConnection(root ComponentFn[Empty], conn net.Conn) {
	defer func() {
		recover()
		s.close()
	}()

	s.assignments = NewPatches(s.onSendTemplatePatch)
	s.conn = conn
	renderable := root(s, Empty{})
	s.track(renderable.ref)
	assignments := renderable.assignments()
	s.initialized = true

	s.deferredAssigns.Wait()
	s.send(s.patchesToJson(assignments))

	for {
		msg, op, err := wsutil.ReadClientData(s.conn)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			if _, ok := err.(wsutil.ClosedError); ok {
				break
			}
		}
		if op == ws.OpText {
			ev := Event{}
			if err := json.Unmarshal(msg, &ev); err != nil {
				log.Println(err)
			}
			s.handleEvent(ev)
		}
	}
}

func (s *socket) templateIsKnown(tmpl *Template) bool {
	return s.knownTemplates.has(tmpl.ID())
}

func (s *socket) markTemplateAsKnown(tmpl *Template) {
	s.knownTemplates.set(tmpl.ID())
}

func (s *socket) send(data []byte) {
	s.sending.Lock()
	if s.conn != nil {
		wsutil.WriteServerText(s.conn, data)
	}
	s.sending.Unlock()
}

func (s *socket) onSendTemplatePatch(ps []*assignment) {
	s.send(s.patchesToJson(ps))
}

func (s *socket) patchesToJson(ps []*assignment) []byte {
	templatesToSend := []json.RawMessage{}
	dataToSend := []json.RawMessage{}
	for i := range ps {
		for _, v := range ps[i].data {
			switch av := v.(type) {
			case AppliedView:
				if !s.templateIsKnown(av.template) {
					t, _ := av.marshalAssignmentJSON()
					templatesToSend = append(templatesToSend, t)
					s.markTemplateAsKnown(av.template)
				}
			case []AppliedView:
				for v := range av {
					e := av[v]
					if !s.templateIsKnown(e.template) {
						t, _ := e.marshalAssignmentJSON()
						templatesToSend = append(templatesToSend, t)
						s.markTemplateAsKnown(e.template)
					}
				}
			}
		}
		d := make([]json.RawMessage, 2)
		d[0] = ps[i].socketid
		d[1], _ = json.Marshal(ps[i].data)
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

func (s *socket) track(sock Socket) {
	s.tracking.Lock()
	id := sock.id()
	r := s.refs[id]
	r.socket = sock
	r.count = r.count + 1
	s.refs[id] = r
	s.check(id)
	s.tracking.Unlock()
}

func (s *socket) untrack(sock Socket) {
	s.tracking.Lock()
	id := sock.id()
	cs := s.refs[id]
	cs.count = cs.count - 1
	s.refs[id] = cs
	s.check(id)
	s.tracking.Unlock()
}

func (s *socket) check(id uintptr) {
	if r, ok := s.refs[id]; ok && r.count < 1 {
		r.socket.dispatch(UnmountEvent, nil)
		delete(s.refs, id)
	}
}

func (s *socket) close() {
	s.dispatch(UnmountEvent, nil)
	if s.conn != nil {
		s.conn.Close()
	}
}

func (s *socket) dispatch(ev string, data EventData) {
	for i, v := range s.refs {
		if v.count < 1 {
			continue
		}
		sr, ok := s.refs[i]
		if ok {
			sr.count = 0
			s.refs[i] = sr
			sr.socket.dispatch(ev, data)
		}
	}
}
