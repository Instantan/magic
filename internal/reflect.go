package internal

import (
	"reflect"
	"strconv"
	"strings"
)

func Traverse(obj any, callback func(value any, path string)) {
	traverseRecursive(reflect.ValueOf(obj), "", callback)
}

func traverseRecursive(value reflect.Value, path string, callback func(value any, path string)) {
	switch value.Kind() {
	case reflect.Ptr:
		originalValue := value.Elem()
		if !originalValue.IsValid() {
			return
		}
		traverseRecursive(originalValue, path, callback)
	case reflect.Interface:
		traverseRecursive(value.Elem(), path, callback)
	case reflect.Struct:
		for _, f := range reflect.VisibleFields(value.Type()) {
			if f.IsExported() {
				name := f.Name
				if v := f.Tag.Get("json"); len(v) > 0 {
					parts := strings.Split(v, ",")
					name = parts[0]
				}
				p := path + "." + name
				if p[0] == '.' {
					p = p[1:]
				}
				traverseRecursive(value.FieldByIndex(f.Index), p, callback)
			}
		}
	case reflect.Slice:
		for i := 0; i < value.Len(); i += 1 {
			p := path + ".[" + strconv.Itoa(i) + "]"
			if p[0] == '.' {
				p = p[1:]
			}
			traverseRecursive(value.Index(i), p, callback)
		}
	case reflect.Map:
		for _, key := range value.MapKeys() {
			p := path + "." + key.String()
			if p[0] == '.' {
				p = p[1:]
			}
			traverseRecursive(value.MapIndex(key), p, callback)
		}
	default:
		zero := reflect.Value{}
		if value != zero {
			callback(value.Interface(), path)
		}
	}
}
