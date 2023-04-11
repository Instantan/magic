package magic

import (
	"encoding/json"
	"sync"
)

type Reactive interface {
	Subscribe(Patchable)
	Unsubscribe(Patchable)
}

type ReactiveValue[T any] interface {
	Reactive
	Get() T
	Set(v T)
}

type value[T any] struct {
	v           *T
	subscribers *map[Patchable]struct{}
}

type atomicValue[T any] struct {
	v           *T
	l           *sync.RWMutex
	subscribers *map[Patchable]struct{}
}

func Value[T any](v T) ReactiveValue[T] {
	return value[T]{
		v:           &v,
		subscribers: &map[Patchable]struct{}{},
	}
}

func AtomicValue[T any](v T) ReactiveValue[T] {
	return atomicValue[T]{
		v:           &v,
		l:           &sync.RWMutex{},
		subscribers: &map[Patchable]struct{}{},
	}
}

func (v value[T]) Set(value T) {
	(*v.v) = value
	v.patch(Rpl, "", value)
}

func (v value[T]) Get() T {
	return *v.v
}

func (v value[T]) Subscribe(patchable Patchable) {
	if (*v.subscribers) == nil {
		(*v.subscribers) = map[Patchable]struct{}{}
	}
	(*v.subscribers)[patchable] = struct{}{}
}

func (v value[T]) Unsubscribe(patchable Patchable) {
	if (*v.subscribers) == nil {
		(*v.subscribers) = map[Patchable]struct{}{}
	}
	delete((*v.subscribers), patchable)
}

func (v value[T]) patch(op Operation, path string, data any) {
	for patchable := range *v.subscribers {
		patchable.Patch(op, path, data)
	}
}

func (v value[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.Get())
}

func (v atomicValue[T]) Set(value T) {
	v.l.Lock()
	(*v.v) = value
	v.l.RUnlock()
	v.patch(Rpl, "", value)
}

func (v atomicValue[T]) Get() T {
	v.l.RLock()
	val := *v.v
	v.l.RUnlock()
	return val
}

func (v atomicValue[T]) Subscribe(patchable Patchable) {
	v.l.Lock()
	(*v.subscribers)[patchable] = struct{}{}
	v.l.Unlock()
}

func (v atomicValue[T]) Unsubscribe(patchable Patchable) {
	v.l.Lock()
	delete((*v.subscribers), patchable)
	v.l.Unlock()
}

func (v atomicValue[T]) patch(op Operation, path string, data any) {
	v.l.RLock()
	for patchable := range *v.subscribers {
		patchable.Patch(op, path, data)
	}
	v.l.RUnlock()
}

func (v atomicValue[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.Get())
}
