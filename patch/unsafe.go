package patch

var tracked = map[uintptr]*Patches{}

func PushPatchOf(of uintptr, op Operation, path string, value any) {
	if tracked[of] == nil {
		tracked[of] = &Patches{}
	}
	tracked[of].PushPatch(op, path, value)
}

func PatchesOf(of uintptr) *Patches {
	return tracked[of]
}
