package main

import (
	"encoding/json"
	"unsafe"

	"github.com/Instantan/magic"
	"github.com/Instantan/magic/patch"
)

func main() {

	l := magic.List[int]{}

	l.Append(1)

	patches := patch.PatchesOf(uintptr(unsafe.Pointer(&l)))

	b, err := json.Marshal(patches)
	if err != nil {
		panic(err)
	}
	println(string(b))
}
