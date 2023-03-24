package magic

import (
	"encoding/json"
)

type Operation byte

type Patchable interface {
	Patch(op Operation, path string, data any)
}

type PatchRedirecter struct {
	Patchable Patchable
	Path      string
}

type PatchReceiver func(op Operation, path string, data any)

const (
	Add Operation = iota
	Del
	Rpl
	Swp
)

type Patch struct {
	Op    Operation
	Path  string
	Value any
}

type Patches []Patch

func (p *Patches) PushPatch(op Operation, path string, value any) {
	*p = append(*p, Patch{Op: op, Path: path, Value: value})
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
	d, err := json.Marshal(p.Value)
	if err != nil {
		return nil, err
	}
	switch p.Op {
	case Add:
		return joinBytesSlicesAndSetLastToCloseBrace([]byte("[0,\""+p.Path+"\","), d), nil
	case Del:
		return joinBytesSlicesAndSetLastToCloseBrace([]byte("[1,\""+p.Path+"\","), d), nil
	case Rpl:
		return joinBytesSlicesAndSetLastToCloseBrace([]byte("[2,\""+p.Path+"\","), d), nil
	case Swp:
		return joinBytesSlicesAndSetLastToCloseBrace([]byte("[3,\""+p.Path+"\","), d), nil
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

func (p Patch) String() string {
	b, err := p.MarshalJSON()
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

func (p *PatchRedirecter) Patch(op Operation, path string, value any) {
	p.Patchable.Patch(op, p.Path+path, value)
}

func (p *PatchReceiver) Patch(op Operation, path string, value any) {
	(*p)(op, path, value)
}
