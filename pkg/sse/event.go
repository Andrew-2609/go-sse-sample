package sse

type Event struct {
	Type EventType `json:"type,omitempty"`
	Data any       `json:"data"`
}

type EventType string

const (
	EventTypeNone EventType = ""
)

func (e *Event) IsEmpty() bool {
	return e.Data == nil
}
