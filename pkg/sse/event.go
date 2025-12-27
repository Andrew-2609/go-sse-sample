package sse

import (
	"log"

	"github.com/google/uuid"
)

type Event struct {
	ID   string    `json:"id,omitempty"`
	Type EventType `json:"type,omitempty"`
	Data any       `json:"data"`
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
		ID:   id,
		Type: eventType,
		Data: data,
	}
}

func (e *Event) IsEmpty() bool {
	return e.Data == nil
}
