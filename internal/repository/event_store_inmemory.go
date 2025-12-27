package repository

import (
	"log"
	"sync"
	"time"

	"github.com/Andrew-2609/go-sse-sample/pkg/sse"
)

type EventStoreInMemory struct {
	mu            sync.Mutex
	events        []sse.Event
	stopRetention chan struct{}
}

var _ sse.EventStore = (*EventStoreInMemory)(nil)

func NewEventStoreInMemory(ttl time.Duration) *EventStoreInMemory {
	store := &EventStoreInMemory{
		mu:            sync.Mutex{},
		events:        make([]sse.Event, 0),
		stopRetention: make(chan struct{}, 1),
	}

	store.startRetention(ttl)

	return store
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

func (e *EventStoreInMemory) startRetention(ttl time.Duration) {
	ticker := time.NewTicker(ttl / 2)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				log.Printf("running events retention at %s\n", time.Now().UTC().Format(time.RFC3339))

				cutoff := time.Now().Add(-ttl)

				e.mu.Lock()

				i := 0
				for _, event := range e.events {
					if event.CreatedAt.After(cutoff) {
						break
					}
					i++
				}

				if i > 0 {
					log.Printf("events retention: removing %d events", i)
					e.events = e.events[i:]
				}

				e.mu.Unlock()
			case <-e.stopRetention:
				log.Println("events retention stopped")
				return
			}
		}

	}()
}

func (e *EventStoreInMemory) StopRetention() {
	e.stopRetention <- struct{}{}
}
