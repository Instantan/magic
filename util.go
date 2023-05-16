package magic

import (
	"encoding/json"
	"fmt"
)

type Empty = struct{}
type Nothing = Empty

// Must panics if the passed error is not nil
func Must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}

func Map[T any, R any](s []T, f func(e T) R) []R {
	ns := make([]R, len(s))
	for i := range s {
		ns[i] = f(s[i])
	}
	return ns
}

func socketid(id uintptr) json.RawMessage {
	v, _ := json.Marshal(fmt.Sprintf("%v", id))
	return v
}
