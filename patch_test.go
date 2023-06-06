package magic

import "testing"

func TestPatchMarshalJSON(t *testing.T) {
	p := Patch{
		op:     PatchOpDEL,
		target: "[0]",
		data:   0,
	}
	t.Log(string(Must(p.MarshalJSON())))
}
