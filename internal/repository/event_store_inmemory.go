package repository

import (
	"github.com/Andrew-2609/go-sse-sample/pkg/sse"
)

type EventStoreInMemory struct {
	events []sse.Event
}

var _ sse.EventStore = (*EventStoreInMemory)(nil)

func NewEventStoreInMemory() *EventStoreInMemory {
	return &EventStoreInMemory{
		events: make([]sse.Event, 0),
	}
}

func (e *EventStoreInMemory) StoreEvent(event sse.Event) {
	e.events = append(e.events, event)
}

func (e *EventStoreInMemory) GetEventsAfterID(id string) []sse.Event {
	for i, event := range e.events {
		if event.ID == id {
			return e.events[i+1:]
		}
	}
	return nil
}
