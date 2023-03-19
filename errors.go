package magic

import "errors"

var ErrUnreachable = errors.New("Unreachable")

func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}
