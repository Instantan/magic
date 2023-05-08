package magic

import (
	"encoding/json"
	"fmt"
)

// Must panics if the passed error is not nil
func Must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}

func socketid(id uintptr) json.RawMessage {
	v, _ := json.Marshal(fmt.Sprintf("%v", id))
	return v
}

type Empty = struct{}
type Nothing = Empty
