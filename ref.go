package magic

import (
	"net/http"
	"sync"
	"unsafe"
)

type ref struct {
	root         *socket
	eventHandler EventHandler
	state        map[string]any
	assigning    sync.Mutex
}

func (s *ref) DispatchEvent(ev string, data any) error {
	if !s.Live() {
		return nil
	}
	return s.root.dispatchEvent(ev, data, s.id())
}

func (s *ref) HandleEvent(evh EventHandler) {
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

func (s *ref) Live() bool {
	return s.root.Live()
}

func (s *ref) Request() *http.Request {
	return s.root.Request()
}

func (s *ref) id() uintptr {
	return uintptr(unsafe.Pointer(s))
}

func (s *ref) clone() Socket {
	return &ref{
		root:  s.root,
		state: map[string]any{},
	}
}

func (s *ref) assign(key string, value any) {
	if isDeferred(value) {
		s.assignDeferred(key, value)
	} else if s.Live() {
		s.assignLive(key, value)
	} else {
		s.assignStatic(key, value)
	}
}

func (s *ref) assignStatic(key string, value any) {
	s.assigning.Lock()
	s.state[key] = value
	s.assigning.Unlock()
}

func (s *ref) assignLive(key string, value any) {
	s.assigning.Lock()
	prev := s.state[key]
	s.assigning.Unlock()
	if prev == value {
		return
	}
	if av, ok := prev.(AppliedView); ok {
		av.ref.untrack(nil)
	} else if avs, ok := prev.([]AppliedView); ok {
		for v := range avs {
			avs[v].ref.untrack(nil)
		}
	}
	s.assigning.Lock()
	s.state[key] = value
	s.assigning.Unlock()
	if av, ok := value.(AppliedView); ok {
		s.track(av.ref)
	} else if avs, ok := value.([]AppliedView); ok {
		for v := range avs {
			s.track(avs[v].ref)
		}
	}
	if s.root != nil && s.root.conn != nil && s.root.assignments != nil && s.root.initialized {
		p := getAssignment()
		p.socketid = socketid(s.id())
		p.data = map[string]any{
			key: value,
		}
		s.root.assignments.append(p)
	}
}

func (s *ref) assignDeferred(key string, value any) {
	s.root.deferredAssigns.Add(1)
	submitTask(func() {
		s.assign(key, resolveDeferred(value))
		s.root.deferredAssigns.Done()
	})
}

func (s *ref) track(sock Socket) {
	s.root.track(sock)
}

func (s *ref) untrack(sock Socket) {
	if sock == nil {
		s.assigning.Lock()
		for _, v := range s.state {
			switch v := v.(type) {
			case AppliedView:
				v.ref.untrack(nil)
			case []AppliedView:
				for i := range v {
					v[i].ref.untrack(nil)
				}
			}
		}
		s.assigning.Unlock()
		s.root.untrack(s)
		return
	}
	s.root.untrack(sock)
}

func (s *ref) dispatch(ev string, data EventData) {
	if s.eventHandler == nil {
		return
	}
	s.eventHandler(ev, data)
}
