package magic

import "unsafe"

type socketref struct {
	root         *socket
	eventHandler EventHandler
	state        map[string]any
}

func (s *socketref) Send(ev string, data any) error {
	return s.root.Send(ev, data)
}

func (s *socketref) HandleEvent(evh EventHandler) {
	s.eventHandler = evh
}

func (s *socketref) Live() bool {
	return s.root.Live()
}

func (s *socketref) Done() <-chan struct{} {
	return s.root.Done()
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

func (s *socketref) assign(key string, value any) {
	prev := s.state[key]
	if prev == value {
		return
	}
	s.state[key] = value
	if s.root != nil && s.root.conn != nil {
		p := getPatch()
		_, refid := s.id()
		p.socketid = socketid(refid)
		p.data = map[string]any{}
		p.data[key] = value
		s.root.patches.append(p)
	}
}
