package magic

import (
	"net"

	"github.com/goccy/go-json"
)

// A line is a live connection to a client
// it is responsible for sending and receiving events

type Line struct {
	state any
	patch *patcher
	conn  net.Conn
}

type event struct {
	Event string          `json:"e"`
	Data  json.RawMessage `json:"d"`
}

// The call to sync, computes a patch between the state of the client and the new state
// ands sends the patch to the client
func (l *Line) Sync() {
	diff, err := l.patch.diff(l.state)
	if err != nil {
		// todo: handle that
		panic(err)
	}
	l.sendEvent(event{
		Event: "p",
		Data:  diff,
	})
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
