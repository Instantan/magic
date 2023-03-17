package magic

type Value[T any] struct {
	value T
}

// Value returns a pointer to the inner value
// Be carefull, manipulating the value directly doesnt result in a change of the clients state
// use is for reading only
func (v *Value[T]) Value() *T {
	return &v.value
}

// Set updates the underlying value
func (v *Value[T]) Set(value T) {
	v.value = value
}
