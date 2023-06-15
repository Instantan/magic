package magic

import (
	"time"
)

type Value interface {
	int |
		int8 |
		int16 |
		int32 |
		int64 |
		uint |
		uint8 |
		uint16 |
		uint32 |
		uint64 |
		uintptr |
		float32 |
		float64 |
		complex64 |
		complex128 |
		string |
		bool |
		time.Time |
		AppliedView |
		[]AppliedView |
		func() int |
		func() int8 |
		func() int16 |
		func() int32 |
		func() int64 |
		func() uint |
		func() uint8 |
		func() uint16 |
		func() uint32 |
		func() uint64 |
		func() uintptr |
		func() float32 |
		func() float64 |
		func() complex64 |
		func() complex128 |
		func() string |
		func() bool |
		func() time.Time |
		func() AppliedView |
		func() []AppliedView |
		<-chan int |
		<-chan int8 |
		<-chan int16 |
		<-chan int32 |
		<-chan int64 |
		<-chan uint |
		<-chan uint8 |
		<-chan uint16 |
		<-chan uint32 |
		<-chan uint64 |
		<-chan uintptr |
		<-chan float32 |
		<-chan float64 |
		<-chan complex64 |
		<-chan complex128 |
		<-chan string |
		<-chan bool |
		<-chan time.Time |
		<-chan AppliedView |
		<-chan []AppliedView
}

// Assign assigns a new value to the given socket
func Assign[T Value](s Socket, key string, value T) {
	s.assign(key, value)
}

func isDeferred(v any) bool {
	switch v.(type) {
	case func() int:
		return true
	case func() int8:
		return true
	case func() int16:
		return true
	case func() int32:
		return true
	case func() int64:
		return true
	case func() uint:
		return true
	case func() uint8:
		return true
	case func() uint16:
		return true
	case func() uintptr:
		return true
	case func() float32:
		return true
	case func() float64:
		return true
	case func() complex64:
		return true
	case func() complex128:
		return true
	case func() string:
		return true
	case func() bool:
		return true
	case func() time.Time:
		return true
	case func() AppliedView:
		return true
	case func() []AppliedView:
		return true
	case <-chan int:
		return true
	case <-chan int8:
		return true
	case <-chan int16:
		return true
	case <-chan int32:
		return true
	case <-chan int64:
		return true
	case <-chan uint:
		return true
	case <-chan uint8:
		return true
	case <-chan uint16:
		return true
	case <-chan uintptr:
		return true
	case <-chan float32:
		return true
	case <-chan float64:
		return true
	case <-chan complex64:
		return true
	case <-chan complex128:
		return true
	case <-chan string:
		return true
	case <-chan bool:
		return true
	case <-chan time.Time:
		return true
	case <-chan AppliedView:
		return true
	case <-chan []AppliedView:
		return true
	default:
		return false
	}
}

func resolveDeferred(v any) any {
	switch d := v.(type) {
	case func() int:
		return d()
	case func() int8:
		return d()
	case func() int16:
		return d()
	case func() int32:
		return d()
	case func() int64:
		return d()
	case func() uint:
		return d()
	case func() uint8:
		return d()
	case func() uint16:
		return d()
	case func() uintptr:
		return d()
	case func() float32:
		return d()
	case func() float64:
		return d()
	case func() complex64:
		return d()
	case func() complex128:
		return d()
	case func() string:
		return d()
	case func() bool:
		return d()
	case func() time.Time:
		return d()
	case func() AppliedView:
		return d()
	case func() []AppliedView:
		return d()
	case <-chan int:
		return <-d
	case <-chan int8:
		return <-d
	case <-chan int16:
		return <-d
	case <-chan int32:
		return <-d
	case <-chan int64:
		return <-d
	case <-chan uint:
		return <-d
	case <-chan uint8:
		return <-d
	case <-chan uint16:
		return <-d
	case <-chan uintptr:
		return <-d
	case <-chan float32:
		return <-d
	case <-chan float64:
		return <-d
	case <-chan complex64:
		return <-d
	case <-chan complex128:
		return <-d
	case <-chan string:
		return <-d
	case <-chan bool:
		return <-d
	case <-chan time.Time:
		return <-d
	case <-chan AppliedView:
		return <-d
	case <-chan []AppliedView:
		return <-d
	default:
		return ""
	}
}
