package magic

import (
	"encoding/json"
	"log"
)

const (
	UnmountEvent         = "m:unmount"
	ClickEvent           = "m:click"
	FocusEvent           = "m:focus"
	ChangeEvent          = "m:change"
	KeydownEvent         = "m:keydown"
	KeypressEvent        = "m:keypress"
	KeyupEvent           = "m:keyup"
	SubmitEvent          = "m:submit"
	DblclickEvent        = "m:dblclick"
	NavigateEvent        = "m:navigate"
	ResetEvent           = "m:reset"
	BlurEvent            = "m:blur"
	OpenFullscreenEvent  = "m:openFullscreen"
	CloseFullscreenEvent = "m:closeFullscreen"
	AnimateEvent         = "m:animate"
	ScrollIntoViewEvent  = "m:scrollIntoView"
	DisconnectEvent      = "m:disconnect"
)

type Event struct {
	Kind    string          `json:"k"`
	Target  uintptr         `json:"t"`
	Payload json.RawMessage `json:"p"`
}

type EventData json.RawMessage

type EventHandler func(ev string, data EventData)
type EventSender func(ev string, data EventData)

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
	Key     string `json:"key"`
	Content string `json:"content"`
}

type KeydownPayload struct {
	Value   string `json:"value"`
	MetaKey bool   `json:"metaKey"`
	CtrlKey bool   `json:"shiftKey"`
	Key     string `json:"key"`
	Content string `json:"content"`
}

type KeypressPayload struct {
	Value   string `json:"value"`
	MetaKey bool   `json:"metaKey"`
	CtrlKey bool   `json:"shiftKey"`
	Key     string `json:"key"`
	Content string `json:"content"`
}

type KeyupPayload struct {
	Value   string `json:"value"`
	MetaKey bool   `json:"metaKey"`
	CtrlKey bool   `json:"shiftKey"`
	Key     string `json:"key"`
	Content string `json:"content"`
}

type SubmitPayload struct {
	Value   string          `json:"value"`
	MetaKey bool            `json:"metaKey"`
	CtrlKey bool            `json:"shiftKey"`
	Form    json.RawMessage `json:"form"`
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
