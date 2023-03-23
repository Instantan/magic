package magic

import "encoding/json"

type Reactive interface {
	Subscribe(patchable Patchable)
	Unsubscribe(Patchable)
}

type Patchable interface {
	Patch(op Operation, path string, data any)
}

type Get[T any] func(su ...suAction) T
type Set[T any] func(v T)

type suAction struct {
	Patchable Patchable
	subscribe bool
}

type signalSubscriptions struct {
	subscribed map[Patchable]struct{}
}

func CreateSignal[T any](init T) (Get[T], Set[T]) {
	value := init
	ssub := signalSubscriptions{
		subscribed: map[Patchable]struct{}{},
	}
	getter := Get[T](func(su ...suAction) T {
		if len(su) > 0 {
			ssub.apply(su)
		}
		return value
	})
	setter := Set[T](func(v T) {
		ssub.patch(Rpl, "", v)
		value = v
	})
	return getter, setter
}

func (g Get[T]) Subscribe(patchable Patchable) {
	g(suAction{
		Patchable: patchable,
		subscribe: true,
	})
}

func (g Get[T]) Unsubscribe(patchable Patchable) {
	g(suAction{
		Patchable: patchable,
		subscribe: false,
	})
}

func (g Get[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(g())
}

func (s *signalSubscriptions) apply(su []suAction) {
	for _, action := range su {
		if action.subscribe {
			s.subscribed[action.Patchable] = struct{}{}
		} else {
			delete(s.subscribed, action.Patchable)
		}
	}
}

func (s *signalSubscriptions) patch(op Operation, path string, data any) {
	for patchable := range s.subscribed {
		patchable.Patch(op, path, data)
	}
}
