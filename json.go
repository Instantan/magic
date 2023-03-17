package magic

import (
	"github.com/goccy/go-json"

	jsonpatch "github.com/evanphx/json-patch/v5"
)

// The heart of magic is the mechanism of the liveview
// instead of sending html it sends a json patch of the pages app state
// the rendering is done by the client with a few primitives like
// show (if)
// range (foreach)
// value
// the primitives are implemented as custom components in html

type patcher struct {
	previous []byte
}

func (p *patcher) diff(data any) ([]byte, error) {
	n, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	patch, err := jsonpatch.CreateMergePatch(p.previous, n)
	if err != nil {
		return nil, err
	}
	p.previous = n
	return patch, nil
}
