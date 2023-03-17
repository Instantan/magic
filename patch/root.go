package patch

import (
	"reflect"
	"unsafe"

	"github.com/viant/xunsafe"
)

type Root struct {
	patches Patches
}

func NewRoot(data any) Root {
	r := Root{}
	r.register(data, reflect.TypeOf(data))
	return r
}

func (r *Root) Patches() Patches {
	return r.patches
}

func (r *Root) PushPatch(op Operation, path string, value any) {
	r.patches.PushPatch(op, path, value)
}

func (r *Root) RegisterParent(path string, _ Patchable) {
}

func (r *Root) register(data any, typ reflect.Type) {
	switch typ.Kind() {
	// case reflect.Pointer
	case reflect.Struct:
		r.registerStruct(data, typ)
	}
}

func (r *Root) registerStruct(data any, typ reflect.Type) {
	xstruct := xunsafe.NewStruct(typ)
	for i := range xstruct.Fields {
		inf := xstruct.Fields[i].Interface(unsafe.Pointer(&data))
		if patchable, ok := inf.(Patchable); ok {
			name := xstruct.Fields[i].Name
			patchable.RegisterParent("/"+name, r)
		} else {
			r.register(data, typ)
		}
	}
}
