package magic

import "unsafe"

type Context interface {
	id() (root uintptr, self uintptr)
	clone() Context
}

type context struct {
	root *context
}

func (ctx *context) id() (root uintptr, self uintptr) {
	self = uintptr(unsafe.Pointer(ctx))
	if ctx.root == nil {
		return 0, self
	}
	return uintptr(unsafe.Pointer(ctx.root)), self
}

func (ctx *context) clone() Context {
	return &context{
		root: ctx.root,
	}
}
