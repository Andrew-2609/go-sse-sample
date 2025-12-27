package sse

import (
	"log"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID        string    `json:"id,omitempty"`
	Type      EventType `json:"event,omitempty"` // having it serialized as "event" is compliant with the EventSource spec
	Data      any       `json:"data"`
	CreatedAt time.Time `json:"-"`
}

type EventType string

const (
	EventTypeNone EventType = ""
)

func NewEvent(eventType EventType, data any) Event {
	id := uuid.New().String()

	uuidV7, err := uuid.NewV7()
	if err != nil {
		log.Printf("error creating v7 UUID for event id: %v. A default UUID will be used instead.\n", err)
	} else {
		id = uuidV7.String()
	}

	return Event{
		ID:        id,
		Type:      eventType,
		Data:      data,
		CreatedAt: time.Now().UTC(),
	}
}

func (e *Event) IsEmpty() bool {
	return e.Data == nil
}
