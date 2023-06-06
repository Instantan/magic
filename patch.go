package magic

import (
	"encoding/json"
	"fmt"
)

type PatchOp byte

const (
	PatchOpDEL = iota
	PatchOpINS
	PatchOpSWP
)

type Patch struct {
	op     PatchOp
	target string
	data   any
}

func (p Patch) MarshalJSON() ([]byte, error) {
	op, err := json.Marshal(p.op)
	if err != nil {
		return []byte{}, err
	}
	idx, err := json.Marshal(p.target)
	if err != nil {
		return []byte{}, err
	}
	b, err := json.Marshal(p.data)
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal([]json.RawMessage{
		op,
		idx,
		b,
	})
}

func NewPatchSliceDEL(index int) Patch {
	return Patch{
		op:     PatchOpDEL,
		target: fmt.Sprintf("[%v]", index),
		data:   nil,
	}
}

func NewPatchSliceSWP(index1, index2 int) Patch {
	return Patch{
		op:     PatchOpSWP,
		target: fmt.Sprintf("[%v]", index1),
		data:   fmt.Sprintf("[%v]", index1),
	}
}

func NewPatchSliceINS(index int, appliedView AppliedView) Patch {
	return Patch{
		op:     PatchOpSWP,
		target: fmt.Sprintf("[%v]", index),
		data:   appliedView,
	}
}
