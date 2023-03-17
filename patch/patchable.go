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

type Patchable interface {
	PushPatch(op Operation, path string, value any)
	RegisterParent(path string, p Patchable)
}

// type Patcher struct {
// 	path   string
// 	parent Patchable
// }

// func (p *Patcher) PushPatch(op Operation, path string, value any) {
// 	if p.parent != nil {
// 		p.parent.PushPatch(op, path, value)
// 	}
// }

// func (p *Patcher) RegisterParent(path string, pa Patchable) {
// 	p.path = path
// 	p.parent = pa
// }
