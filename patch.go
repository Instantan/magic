package magic

import (
	"encoding/json"
)

type Operation byte

type PatchRedirecter struct {
	Patchable Patchable
	Path      string
}

const (
	Add Operation = iota
	Del
	Rpl
	Swp
)

type Patch struct {
	op    Operation
	path  string
	value any
}

type ChangeTracker interface {
	PushPatch(root uintptr, op Operation, path string, value any)
}

type Patches []Patch

func (p *Patches) PushPatch(op Operation, path string, value any) {
	*p = append(*p, Patch{op: op, path: path, value: value})
}

func (p Patches) MarshalJSON() ([]byte, error) {
	if len(p) == 0 {
		return []byte("[]"), nil
	}
	data := []byte{0}
	data[0] = 91
	for i := range p {
		d, err := p[i].MarshalJSON()
		if err != nil {
			return nil, err
		}
		data = append(data, d...)
		data = append(data, 44)
	}
	data[len(data)-1] = 93
	return data, nil
}

func (p *Patch) MarshalJSON() ([]byte, error) {
	d, err := json.Marshal(p.value)
	if err != nil {
		return nil, err
	}
	switch p.op {
	case Add:
		return joinBytesSlicesAndSetLastToCloseBrace([]byte("[0,\""+p.path+"\","), d), nil
	case Del:
		return joinBytesSlicesAndSetLastToCloseBrace([]byte("[1,\""+p.path+"\","), d), nil
	case Rpl:
		return joinBytesSlicesAndSetLastToCloseBrace([]byte("[2,\""+p.path+"\","), d), nil
	case Swp:
		return joinBytesSlicesAndSetLastToCloseBrace([]byte("[3,\""+p.path+"\","), d), nil
	default:
		return nil, ErrUnreachable
	}
}

func (p Patches) String() string {
	b, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(b)
}

func joinBytesSlicesAndSetLastToCloseBrace(s1, s2 []byte) []byte {
	n := len(s1)
	n += len(s2)
	b, i := make([]byte, n+1), 0
	i += copy(b[i:], s1)
	i += copy(b[i:], s2)
	b[n] = 93
	return b
}

func (p *PatchRedirecter) Patch(patch Patch) {
	patch.path = p.Path + patch.path
	p.Patchable.Patch(patch)
}
