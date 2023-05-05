package magic

import (
	"reflect"
	"unsafe"
)

func unsafeBytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func unsafeStringToBytes(s string) (b []byte) {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	return b
}
