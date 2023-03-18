package patch

import (
	"reflect"
)

type Root struct {
	patches Patches
}

func NewRoot(data any) Root {
	r := Root{}
	r.register(data, reflect.ValueOf(data))
	return r
}

func (r *Root) Patches() Patches {
	return r.patches
}

func (r *Root) PushPatch(op Operation, path string, value any) {
	r.patches.PushPatch(op, path, value)
}

func (r *Root) RegisterParent(path string, _ Tracked) {
}

func (r *Root) register(data any, typ reflect.Value) {
	switch typ.Kind() {
	// case reflect.Pointer
	case reflect.Struct:
		r.registerStruct(data, typ)
	}
}

func (r *Root) registerStruct(data any, v reflect.Value) {
	n := v.NumField()
	t := v.Type()
	for i := 0; i < n; i++ {
		f := v.Field(i)
		if t.Field(i).IsExported() {
			inf := f.Interface()
			if patchable, ok := inf.(Tracked); ok {
				name := t.Field(i).Name
			} else {
				r.register(data, f)
			}
		}
	}
}

// Register registers the path of a
func Register(parent uintptr, path string, value any) {
	"name/"
}
