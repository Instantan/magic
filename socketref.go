package magic

import (
	"unsafe"
)

type socketref struct {
	root         *socket
	eventHandler EventHandler
	state        map[string]any
}

func (s *socketref) DispatchEvent(ev string, data any) error {
	return s.root.DispatchEvent(ev, data)
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
	switch s.Live() {
	case true:
		s.assignLive(key, value)
	case false:
		s.assignStatic(key, value)
	}
}

func (s *socketref) assignStatic(key string, value any) {
	s.state[key] = value
}

func (s *socketref) assignLive(key string, value any) {
	prev := s.state[key]
	if prev == value {
		return
	}
	s.state[key] = value
	if av, ok := value.(AppliedView); ok {
		s.track(av.socketref)
	} else if avs, ok := value.([]AppliedView); ok {
		for v := range avs {
			s.track(avs[v].socketref)
		}
	}
	if av, ok := prev.(AppliedView); ok {
		av.socketref.untrack(nil)
		s.untrack(av.socketref)
	} else if avs, ok := value.([]AppliedView); ok {
		for v := range avs {
			avs[v].socketref.untrack(nil)
			s.untrack(avs[v].socketref)
		}
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
	if sock == nil {
		for _, v := range s.state {
			if v, ok := v.(AppliedView); ok && v.socketref != s {
				s.untrack(v.socketref)
				v.socketref.untrack(nil)
			}
		}
		return
	}
	s.root.untrack(sock)
}

func (s *socketref) dispatch(ev string, data EventData) {
	if s.eventHandler == nil {
		return
	}
	s.eventHandler(ev, data)
}
