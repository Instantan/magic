package magic

import (
	"testing"

	"github.com/Instantan/magic/patch"
)

type testStruct1 struct {
	L1 List[int]
}

func TestRegisterStruct(t *testing.T) {
	t1 := testStruct1{
		L1: List[int]{},
	}

	root := patch.NewRoot(t1)

	t1.L1.Append(1)

	t.Log(t1)
	t.Log(root.Patches())
}
