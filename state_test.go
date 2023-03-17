package magic

import (
	"testing"

	"github.com/Instantan/magic/patch"
)

type testStruct1 struct {
	l1 *List[int]
}

func TestRegisterStruct(t *testing.T) {
	t1 := testStruct1{
		l1: &List[int]{},
	}

	root := patch.NewRoot(t1)

	t.Log(root)
}
