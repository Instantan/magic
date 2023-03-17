package magic

import (
	"errors"

	"github.com/goccy/go-json"
)

type op byte

var ErrUnreachable = errors.New("Unreachable")

const (
	add op = iota
	del
	rpl
	swp
)

type patchable interface {
	push(op op, path string, value any) patchable
}

type patch struct {
	op    op
	path  string
	value any
}

type patchesref struct {
	path   string
	parent patchable
}

type patches []patch

// func (p *patches) push(op op, path string, value any) *patches {
// 	*p = append(*p, patch{op: op, path: path, value: value})
// 	return p
// }

func (p *patchesref) push(op op, path string, value any) *patchesref {
	p.parent.push(op, p.path+path, value)
	return p
}

func (p *patches) clear() *patches {
	*p = []patch{}
	return p
}

func (p patches) MarshalJSON() ([]byte, error) {
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

func (p *patch) MarshalJSON() ([]byte, error) {
	d, err := json.Marshal(p.value)
	if err != nil {
		return nil, err
	}
	switch p.op {
	case add:
		return joinBytesSlicesAndSetLastToCloseBrace([]byte("[0,\""+p.path+"\","), d), nil
	case del:
		return joinBytesSlicesAndSetLastToCloseBrace([]byte("[1,\""+p.path+"\","), d), nil
	case rpl:
		return joinBytesSlicesAndSetLastToCloseBrace([]byte("[2,\""+p.path+"\","), d), nil
	case swp:
		return joinBytesSlicesAndSetLastToCloseBrace([]byte("[3,\""+p.path+"\","), d), nil
	default:
		return nil, ErrUnreachable
	}
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
