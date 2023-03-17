package magic

import "testing"

func TestPatchPush(t *testing.T) {
	ps := patchesref{}
	ps.push(add, "/", 1).push(add, "/", 1)
	t.Log(ps)
	// ps.clear()
	t.Log(ps)
}

func TestMarshalPatch(t *testing.T) {
	p := patch{0, "", nil}
	data, err := p.MarshalJSON()
	if err != nil {
		t.Fail()
		t.Log(err)
	}
	t.Log(string(data))
}

// func TestMarshalPatches(t *testing.T) {
// 	ps := patches{}
// 	ps.push(add, "/", 1).push(add, "/", 1)
// 	data, err := ps.MarshalJSON()
// 	if err != nil {
// 		t.Fail()
// 		t.Log(err)
// 	}
// 	t.Log(string(data))
// }
