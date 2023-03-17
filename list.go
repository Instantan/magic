package magic

import "strconv"

type List[T any] struct {
	data    []T
	patches patchesref
}

// Value returns a pointer to the inner value (a slice)
// Be carefull, manipulating the slice directly doesnt result in a change of the clients state
// use is for reading only
func (l *List[T]) Value() *[]T {
	return &l.data
}

// Get returns the element at the given index
func (l *List[T]) Get(index int) T {
	return l.data[index]
}

// First returns the first element of the list
func (l *List[T]) First() T {
	return l.data[0]
}

// First returns the last element of the list
func (l *List[T]) Last() T {
	return l.data[len(l.data)-1]
}

// Len returns the length of the list
func (l *List[T]) Len() int {
	return len(l.data)
}

// Append adds a element to the end of the list
func (l *List[T]) Append(v T) *List[T] {
	l.data = append(l.data, v)
	l.patches.push(add, "/[]", v)
	return l
}

// Prepend adds a element to the beginning of the list
func (l *List[T]) Prepend(v T) *List[T] {
	l.data = append([]T{v}, l.data...)
	l.patches.push(add, "/[0]", v)
	return l
}

// Shift removes the first element of the list
func (l *List[T]) Shift() *List[T] {
	if len(l.data) > 0 {
		l.data = l.data[1:]
		l.patches.push(del, "/[0]", nil)
	}
	return l
}

// Pop removes the last element of the list
func (l *List[T]) Pop() *List[T] {
	if len(l.data) > 0 {
		l.data = l.data[:len(l.data)-1]
		l.patches.push(del, "/[]", nil)
	}
	return l
}

// Remove removes the element at the given index from the list
func (l *List[T]) Remove(i int) *List[T] {
	if len(l.data) > i {
		l.data = append(l.data[:i], l.data[i+1:]...)
		l.patches.push(del, "/["+strconv.Itoa(i)+"]", nil)
	}
	return l
}

// Swap switches the value at the position i with the value at the position y
func (l *List[T]) Swap(i, y int) *List[T] {
	l.data[i], l.data[y] = l.data[y], l.data[i]
	l.patches.push(swp, "/["+strconv.Itoa(i)+"]", "/["+strconv.Itoa(i)+"]")
	return l
}

// Set sets the value at the given position
func (l *List[T]) Set(i int, value T) *List[T] {
	l.data[i] = value
	l.patches.push(rpl, "/["+strconv.Itoa(i)+"]", value)
	return l
}

// Nil sets the slice to nil
func (l *List[T]) Nil() *List[T] {
	l.data = nil
	l.patches.push(rpl, "/", nil)
	return l
}

// SetSlice completly changes the underlying slice
func (l *List[T]) SetSlice(slc []T) *List[T] {
	l.data = slc
	l.patches.push(rpl, "/", slc)
	return l
}
