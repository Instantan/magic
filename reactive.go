package magic

import (
	"encoding/json"
	"sync/atomic"
)

type Reactive interface {
	Subscribe(Patchable)
	Unsubscribe(Patchable)
}

type Get[T any] func(s ...s) T
type Set[T any] func(v T)

type ComputedFn[T any] func(s ...s) T

type s struct {
	patchable Patchable
	subscribe bool
}

type signalsubscriptions map[Patchable]struct{}

// Signal is a primitive reactive state
// its only safe to use in a non concurrent scenario
// For a concurrency safe signal use AtomicSignal
func Signal[T any](value T) (Get[T], Set[T]) {
	ssb := signalsubscriptions{}
	getter := Get[T](func(s ...s) T {
		if len(s) > 0 {
			ssb.apply(s)
		}
		return value
	})
	setter := Set[T](func(v T) {
		value = v
		ssb.patch(Rpl, "", v)
	})
	return getter, setter
}

func AtomicSignal[T any](value T) (Get[T], Set[T]) {
	avalue := atomic.Value{}
	avalue.Store(value)
	ssb := signalsubscriptions{}
	getter := Get[T](func(s ...s) T {
		if len(s) > 0 {
			ssb.apply(s)
		}
		return avalue.Load().(T)
	})
	setter := Set[T](func(v T) {
		avalue.Store(v)
		ssb.patch(Rpl, "", v)
	})
	return getter, setter
}

func Computed[T any](fn func() T, deps ...Reactive) Get[T] {
	// TODO: maybe there is a way to do this without introducing a state?
	ssb := signalsubscriptions{}
	pr := PatchReceiver(func(op Operation, path string, data any) {
		ssb.patch(Rpl, "", fn())
	})
	for i := range deps {
		deps[i].Subscribe(&pr)
	}
	return Get[T](func(s ...s) T {
		if len(s) > 0 {
			ssb.apply(s)
		}
		return fn()
	})
}

func Effect(fn func(), deps ...Reactive) {
	pr := PatchReceiver(func(op Operation, path string, data any) {
		fn()
	})
	for i := range deps {
		deps[i].Subscribe(&pr)
	}
}

func (g Get[T]) Subscribe(patchable Patchable) {
	g(s{
		patchable: patchable,
		subscribe: true,
	})
}

func (g Get[T]) Unsubscribe(patchable Patchable) {
	g(s{
		patchable: patchable,
		subscribe: false,
	})
}

func (g Get[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(g())
}

func (ssub *signalsubscriptions) apply(s []s) {
	for _, action := range s {
		if action.subscribe {
			(*ssub)[action.patchable] = struct{}{}
		} else {
			delete((*ssub), action.patchable)
		}
	}
}

func (ssub *signalsubscriptions) patch(op Operation, path string, data any) {
	for patchable := range *ssub {
		patchable.Patch(op, path, data)
	}
}

func (c ComputedFn[T]) Patch(op Operation, path string, data any) {
	c()
}
