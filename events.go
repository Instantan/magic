package magic

import (
	"encoding/json"
	"log"
)

const (
	ClickEvent    = "click"
	FocusEvent    = "focus"
	ChangeEvent   = "change"
	KeydownEvent  = "keydown"
	KeypressEvent = "keypress"
	KeyupEvent    = "keyup"
	SubmitEvent   = "submit"
	DblclickEvent = "dblclick"
)

type Event struct {
	Kind    string          `json:"kind"`
	Payload json.RawMessage `json:"payload"`
}

type EventData json.RawMessage

type EventHandler func(ev string, data any)
type EventSender func(ev string, data any)

type ClickPayload struct {
	Value   string `json:"value"`
	MetaKey bool   `json:"metaKey"`
	CtrlKey bool   `json:"shiftKey"`
}

type FocusPayload struct {
	Value   string `json:"value"`
	MetaKey bool   `json:"metaKey"`
	CtrlKey bool   `json:"shiftKey"`
}

type ChangePayload struct {
	Value   string `json:"value"`
	MetaKey bool   `json:"metaKey"`
	CtrlKey bool   `json:"shiftKey"`
}

type KeydownPayload struct {
	Value   string `json:"value"`
	MetaKey bool   `json:"metaKey"`
	CtrlKey bool   `json:"shiftKey"`
}

type KeypressPayload struct {
	Value   string `json:"value"`
	MetaKey bool   `json:"metaKey"`
	CtrlKey bool   `json:"shiftKey"`
}

type KeyupPayload struct {
	Value   string `json:"value"`
	MetaKey bool   `json:"metaKey"`
	CtrlKey bool   `json:"shiftKey"`
}

type SubmitPayload struct {
	Value   string `json:"value"`
	MetaKey bool   `json:"metaKey"`
	CtrlKey bool   `json:"shiftKey"`
}

type DblclickPayload struct {
	Value   string `json:"value"`
	MetaKey bool   `json:"metaKey"`
	CtrlKey bool   `json:"shiftKey"`
}

func (ev EventData) Click() ClickPayload {
	return parseJsonEventData[ClickPayload](ev)
}

func (ev EventData) Focus() FocusPayload {
	return parseJsonEventData[FocusPayload](ev)
}

func (ev EventData) Change() ChangePayload {
	return parseJsonEventData[ChangePayload](ev)
}

func (ev EventData) Keydown() KeydownPayload {
	return parseJsonEventData[KeydownPayload](ev)
}

func (ev EventData) Keypress() KeypressPayload {
	return parseJsonEventData[KeypressPayload](ev)
}

func (ev EventData) Keyup() KeyupPayload {
	return parseJsonEventData[KeyupPayload](ev)
}

func (ev EventData) Submit() SubmitPayload {
	return parseJsonEventData[SubmitPayload](ev)
}

func (ev EventData) Dblclick() DblclickPayload {
	return parseJsonEventData[DblclickPayload](ev)
}

func parseJsonEventData[T any](data EventData) T {
	t := new(T)
	if err := json.Unmarshal(data, &t); err != nil {
		log.Println(err)
	}
	return *t
}
