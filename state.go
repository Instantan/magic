package magic

import (
	"encoding/json"
	"io"
	"log"
	"time"

	jsonpatch "github.com/evanphx/json-patch/v5"
)

type State[T any] struct {
	Data     T
	previous []byte
	debounce *time.Ticker
	writer   io.Writer
}

// sync computes the patch
func (s *State[T]) Sync() {
	if s.debounce == nil {
		s.diffAndSync()
		return
	}
}

func (s *State[T]) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(s.Data)
	if err != nil {
		return nil, err
	}
	return b, err
}

func (s *State[T]) diffAndSync() {
	n, err := json.Marshal(s.Data)
	if err != nil {
		log.Print(err)
		return
	}
	patch, err := jsonpatch.CreateMergePatch(s.previous, n)
	if err != nil {
		log.Print(err)
		return
	}
	s.previous = n
	s.writer.Write(patch)
}
