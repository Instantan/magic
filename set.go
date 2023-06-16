package magic

import "sync"

type set[T comparable] struct {
	m map[T]struct{}
	l *sync.Mutex
}

func newSet[T comparable]() set[T] {
	return set[T]{
		m: make(map[T]struct{}),
		l: &sync.Mutex{},
	}
}

func (s set[T]) has(v T) bool {
	s.l.Lock()
	_, b := s.m[v]
	s.l.Unlock()
	return b
}

func (s set[T]) set(v T) {
	s.l.Lock()
	s.m[v] = struct{}{}
	s.l.Unlock()
}
