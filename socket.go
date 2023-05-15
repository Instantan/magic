package magic

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

type Socket interface {
	Live() bool
	DispatchEvent(ev string, data any) error
	HandleEvent(EventHandler)

	id() uintptr
	clone() Socket
	assign(key string, value any)

	track(Socket)
	untrack(Socket)
	dispatch(ev string, data EventData)
}

type socket struct {
	socketrefs     map[uintptr]Socket
	socketrefsRefs map[uintptr]uint
	knownTemplates Set[int]

	conn    net.Conn
	patches *patches
}

func (s *socket) Live() bool {
	return s.patches != nil
}

func (s *socket) HandleEvent(evh EventHandler) {
	// we need to dispatch the event to the right ref
}

func (s *socket) handleEvent(ev Event) {
	if ev.Target == 0 {
		s.dispatch(ev.Kind, EventData(ev.Payload))
		return
	}
	sref := s.socketrefs[ev.Target]
	if sref == nil {
		return
	}
	sref.dispatch(ev.Kind, EventData(ev.Payload))
}

func (s *socket) DispatchEvent(ev string, data any) error {
	values, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_ = values
	return nil
}

func (s *socket) id() uintptr {
	return 0
}

func (s *socket) clone() Socket {
	return &socketref{
		root:  s,
		state: map[string]any{},
	}
}

func (s *socket) assign(key string, value any) {
}

func (s *socket) establishConnection(root ComponentFn[Empty], conn net.Conn) {
	defer func() {
		recover()
		s.dispatch(UnmountEvent, nil)
		s.conn.Close()
		s.conn = nil
		s.patches = nil
		s.knownTemplates = Set[int]{}
	}()

	s.patches = NewPatches(s.onSendTemplatePatch)
	renderable := root(s, Empty{})
	s.track(renderable.socketref)
	patches := renderable.Patch()
	s.conn = conn
	data := s.patchesToJson(patches)
	wsutil.WriteServerText(s.conn, data)

	for {
		msg, op, err := wsutil.ReadClientData(s.conn)
		if err != nil {
			if errors.Is(err, io.EOF) {
				continue
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
	return s.knownTemplates.Has(tmpl.ID())
}

func (s *socket) markTemplateAsKnown(tmpl *Template) {
	s.knownTemplates.Set(tmpl.ID())
}

func (s *socket) send(data []byte) {
	s.dumpStore()
	wsutil.WriteServerText(s.conn, data)
}

func (s *socket) onSendTemplatePatch(ps []*patch) {
	s.send(s.patchesToJson(ps))
}

func (s *socket) patchesToJson(ps []*patch) []byte {
	templatesToSend := []json.RawMessage{}
	dataToSend := []json.RawMessage{}
	for i := range ps {
		for _, v := range ps[i].data {
			switch av := v.(type) {
			case AppliedView:
				if !s.templateIsKnown(av.template) {
					m := make([]json.RawMessage, 2)
					m[0], _ = json.Marshal(av.template.ID())
					m[1], _ = json.Marshal(av.template.String())
					t, _ := json.Marshal(m)
					templatesToSend = append(templatesToSend, t)
					s.markTemplateAsKnown(av.template)
				}
			case []AppliedView:
				for v := range av {
					e := av[v]
					if !s.templateIsKnown(e.template) {
						m := make([]json.RawMessage, 2)
						m[0], _ = json.Marshal(e.template.ID())
						m[1], _ = json.Marshal(e.template.String())
						t, _ := json.Marshal(m)
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
	id := sock.id()
	s.socketrefs[id] = sock
	s.socketrefsRefs[id]++
}

func (s *socket) untrack(sock Socket) {
	id := sock.id()
	s.socketrefsRefs[id]--
	if s.socketrefsRefs[id] < 1 {
		if sr := s.socketrefs[id]; sr != nil {
			sr.dispatch(UnmountEvent, nil)
		}
		delete(s.socketrefsRefs, id)
		delete(s.socketrefs, id)
	}
}

func (s *socket) dispatch(ev string, data EventData) {
	for i, v := range s.socketrefsRefs {
		if v < 1 {
			continue
		}
		sr := s.socketrefs[i]
		if sr != nil {
			s.socketrefsRefs[i] = 0
			s.dispatch(ev, data)
		}
	}
}

func (s *socket) dumpStore() {
	log.Printf("Store: %v\n", len(s.socketrefs))
	// for id, s := range s.socketrefs {
	// 	log.Printf("\t%v\n", id)
	// 	for key, value := range s.(*socketref).state {
	// 		log.Printf("\t\t%v - %v\n", key, value)
	// 	}
	// }
}
