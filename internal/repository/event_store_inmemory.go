package repository

import (
	"sync"

	"github.com/Andrew-2609/go-sse-sample/pkg/sse"
)

type EventStoreInMemory struct {
	mu     sync.Mutex
	events []sse.Event
}

var _ sse.EventStore = (*EventStoreInMemory)(nil)

func NewEventStoreInMemory() *EventStoreInMemory {
	return &EventStoreInMemory{
		mu:     sync.Mutex{},
		events: make([]sse.Event, 0),
	}
}

func (e *EventStoreInMemory) StoreEvent(event sse.Event) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.events = append(e.events, event)
}

func (e *EventStoreInMemory) GetEventsAfterID(id string) []sse.Event {
	e.mu.Lock()
	defer e.mu.Unlock()
	for i, event := range e.events {
		if event.ID == id {
			return e.events[i+1:]
		}
	}
	return nil
}
