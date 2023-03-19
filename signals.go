package magic

type Get[T any] func() T
type Set[T any] func(v T)

func CreateSignal[T any](init T) (Get[T], Set[T]) {
	value := init
	getter := func() T {
		return value
	}
	setter := func(v T) {
		value = v
	}
	return getter, setter
}
