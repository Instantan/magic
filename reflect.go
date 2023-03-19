package magic

import (
	"reflect"
)

func Traverse(obj any, callback func(value any, pointer uintptr)) {
	traverseRecursive(reflect.ValueOf(obj), callback)
}

func traverseRecursive(value reflect.Value, callback func(value any, pointer uintptr)) {
	switch value.Kind() {
	case reflect.Ptr:
		originalValue := value.Elem()
		if !originalValue.IsValid() {
			return
		}
		traverseRecursive(originalValue, callback)
	case reflect.Interface:
		traverseRecursive(value.Elem(), callback)
	case reflect.Struct:
		for _, f := range reflect.VisibleFields(value.Type()) {
			if f.IsExported() {
				traverseRecursive(value.FieldByIndex(f.Index), callback)
			}
		}
	case reflect.Slice:
		for i := 0; i < value.Len(); i += 1 {
			traverseRecursive(value.Index(i), callback)
		}
	case reflect.Map:
		for _, key := range value.MapKeys() {
			traverseRecursive(value.MapIndex(key), callback)
		}
	default:
		callback(value.Interface(), uintptr(value.Addr().UnsafePointer()))
	}
}
