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
		string |
		bool |
		time.Time |
		AppliedView |
		[]AppliedView
}

// Assign assigns a new value to the given socket
func Assign[T Value](s Socket, key string, value T) {
	s.assign(key, value)
}
