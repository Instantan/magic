package patch

import (
	"errors"
)

type Operation byte

const (
	Add Operation = iota
	Del
	Rpl
	Swp
)

var ErrUnreachable = errors.New("Unreachable")

type Tracked interface {
	PushPatch(op Operation, path string, value any)
}
