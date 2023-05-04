package magic

import (
	"unsafe"
)

type socketref struct {
	root         *socket
	eventHandler EventHandler
	state        map[string]any
}

func (s *socketref) Send(ev string, data any) error {
	return s.root.Send(ev, data)
}

func (s *socketref) HandleEvent(evh EventHandler) {
	if s.eventHandler == nil {
		s.eventHandler = evh
		return
	}
	old := s.eventHandler
	s.eventHandler = func(ev string, data EventData) {
		evh(ev, data)
		old(ev, data)
	}
}

func (s *socketref) Live() bool {
	return s.root.Live()
}

func (s *socketref) id() uintptr {
	return uintptr(unsafe.Pointer(s))
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
	if av, ok := value.(AppliedView); ok {
		s.track(av.socketref)
	}
	if av, ok := prev.(AppliedView); ok {
		s.untrack(av.socketref)
	}
	if s.root != nil && s.root.conn != nil && s.root.patches != nil {
		p := getPatch()
		p.socketid = socketid(s.id())
		p.data = map[string]any{}
		p.data[key] = value
		s.root.patches.append(p)
	}
}

func (s *socketref) track(sock Socket) {
	s.root.track(sock)
}

func (s *socketref) untrack(sock Socket) {
	s.root.untrack(sock)
}

func (s *socketref) dispatch(ev string, data EventData) {
	if s.eventHandler == nil {
		return
	}
	s.eventHandler(ev, data)
}
