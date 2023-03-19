package magic

import (
	"net"

	"github.com/goccy/go-json"
)

// A line is a live connection to a client
// it is responsible for sending and receiving events

type Line struct {
	state any
	conn  net.Conn
}

type event struct {
	Event string          `json:"e"`
	Data  json.RawMessage `json:"d"`
}

// The dispatch event function can be used to send arbitray data to the client
func (l *Line) DispatchEvent(ev string, data any) {
	db, err := json.Marshal(data)
	if err != nil {
		// todo: handle that
		panic(err)
	}
	l.sendEvent(event{
		Event: ev,
		Data:  db,
	})
}

func (l *Line) sendEvent(ev event) {
	b, err := json.Marshal(ev)
	if err != nil {
		// todo: handle that
		panic(err)
	}
	_, err = l.conn.Write(b)
	if err != nil {
		// todo: handle that
		panic(err)
	}
}

func (ev event) MarshalJSON() ([]byte, error) {
	d, err := json.Marshal(ev.Data)
	if err != nil {
		return nil, err
	}
	return []byte("[" + ev.Event + "," + string(d) + "]"), nil
}

func Append[T any](s *int, slc []T, values ...T) []T {

	return append(slc, values...)
}

type Ref[T any] struct {
	ref *T
}

func RefOf[T any](v *T) Ref[T] {
	return Ref[T]{
		ref: v,
	}
}

func Set[T any](s *int, v T) T {

	Append(s, []int{}, 1)
	return v
}
