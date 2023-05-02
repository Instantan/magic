package magic

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"unsafe"

	"github.com/gobwas/ws/wsutil"
)

type Socket interface {
	Live() bool
	Send(ev string, data any) error
	HandleEvent(EventHandler)

	id() (root uintptr, self uintptr)
	clone() Socket
	assign(key string, value any)
}

type socket struct {
	conn           net.Conn
	knownTemplates Set[int]
	patches        *patches
}

func (s *socket) Live() bool {
	return s.patches != nil
}

func (s *socket) HandleEvent(evh EventHandler) {
	// we need to dispatch the event to the right ref
}

func (s *socket) Send(ev string, data any) error {
	values, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_ = values
	return nil
}

func (s *socket) id() (root uintptr, self uintptr) {
	self = uintptr(unsafe.Pointer(s))
	return self, 0
}

func (s *socket) clone() Socket {
	return &socketref{
		root:  s,
		state: map[string]any{},
	}
}

func (s *socket) assign(key string, value any) {
}

func (s *socket) templateIsKnown(tmpl *Template) bool {
	return s.knownTemplates.Has(tmpl.ID())
}

func (s *socket) markTemplateAsKnown(tmpl *Template) {
	s.knownTemplates.Set(tmpl.ID())
}

func (s *socket) send(data []byte) {
	wsutil.WriteServerText(s.conn, data)
}

func (s *socket) onSendTemplatePatch(ps []*patch) {
	data := s.patchesToJson(ps)
	s.send(data)
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

func socketid(id1, id2 uintptr) json.RawMessage {
	v, _ := json.Marshal(fmt.Sprintf("%v", id2))
	return v
}
