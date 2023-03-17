package magic

import (
	"reflect"
)

type state struct {
}

func newLiveState(data any) {
	typ := reflect.TypeOf(data)
	switch typ.Kind() {
	case reflect.Struct:

	}
	// fooType := xunsafe.TypeOf(data)
}
