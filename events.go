package magic

import "encoding/json"

type Event struct {
	Name  string
	Value any
}

type Events []Event

func (e *Events) PushEvent(name string, value any) {
	*e = append(*e, Event{Name: name, Value: value})
}

func (e *Event) MarshalJSON() ([]byte, error) {
	d, err := json.Marshal(e.Value)
	if err != nil {
		return nil, err
	}
	return joinBytesSlicesAndSetLastToCloseBrace([]byte("[\""+e.Name+"\",\""), d), nil
}

func (e Event) String() string {
	b, err := e.MarshalJSON()
	if err != nil {
		panic(err)
	}
	return string(b)
}

func (e Events) String() string {
	b, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(b)
}
