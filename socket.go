package magic

import (
	"encoding/json"
	"net"
	"unsafe"
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
	conn net.Conn
}

type socketref struct {
	root         *socket
	eventHandler EventHandler
	state        map[string]any
}

func (s *socket) Live() bool {
	return s.conn != nil
}

func (s *socketref) Live() bool {
	return s.root.Live()
}

func (s *socket) HandleEvent(evh EventHandler) {
	// we need to dispatch the event to the right ref
}

func (s *socketref) HandleEvent(evh EventHandler) {
	s.eventHandler = evh
}

func (s *socket) Send(ev string, data any) error {
	values, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_ = values
	return nil
}

func (s *socketref) Send(ev string, data any) error {
	return s.root.Send(ev, data)
}

func (s *socket) id() (root uintptr, self uintptr) {
	self = uintptr(unsafe.Pointer(s))
	return self, 0
}

func (s *socketref) id() (root uintptr, self uintptr) {
	self = uintptr(unsafe.Pointer(s))
	if s.root == nil {
		return 0, self
	}
	return uintptr(unsafe.Pointer(s.root)), self
}

func (s *socketref) clone() Socket {
	return &socketref{
		root:  s.root,
		state: map[string]any{},
	}
}

func (s *socket) clone() Socket {
	return &socketref{
		root:  s,
		state: map[string]any{},
	}
}

func (s *socket) assign(key string, value any) {
}

func (s *socketref) assign(key string, value any) {
	prev := s.state[key]
	// the socket has the same state
	// we dont have to push anything
	if prev == value {
		return
	}
	s.state[key] = value
	if s.Live() {
		// here we need to produce a patch
	}
	// switch v := value.(type) {
	// case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr, float32, float64, bool:
	// 	return
	// case AppliedView:
	// 	return
	// default:
	// }
	// based on the
}
