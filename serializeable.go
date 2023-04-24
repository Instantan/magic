package magic

import (
	"time"
)

type Serializeable interface {
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
		func() string |
		func() bool |
		func() time.Time |
		func() AppliedView |
		func() []AppliedView |
		func(Socket) int |
		func(Socket) int8 |
		func(Socket) int16 |
		func(Socket) int32 |
		func(Socket) int64 |
		func(Socket) uint |
		func(Socket) uint8 |
		func(Socket) uint16 |
		func(Socket) uint32 |
		func(Socket) uint64 |
		func(Socket) uintptr |
		func(Socket) float32 |
		func(Socket) float64 |
		func(Socket) string |
		func(Socket) bool |
		func(Socket) time.Time |
		func(Socket) AppliedView |
		func(Socket) []AppliedView
}

// Assign assigns a new value to the given socket
func Assign[T Serializeable](s Socket, key string, value T) {
	s.assign(key, value)
}
