package magic

import "sync"

type Set[T comparable] struct {
	m map[T]struct{}
	l *sync.Mutex
}

func NewSet[T comparable]() Set[T] {
	return Set[T]{
		m: make(map[T]struct{}),
		l: &sync.Mutex{},
	}
}

func (s Set[T]) Has(v T) bool {
	s.l.Lock()
	_, b := s.m[v]
	s.l.Unlock()
	return b
}

func (s Set[T]) Set(v T) {
	s.l.Lock()
	s.m[v] = struct{}{}
	s.l.Unlock()
}
