package sse

type EventStore interface {
	// StoreEvent stores the event in the event store.
	// It's not guaranteed that the event will be stored successfully.
	//
	// In case of error, the event will be broadcasted anyway. The error must be handled in the store implementation.
	StoreEvent(event Event)

	// GetEventsAfterID returns the events after the given id.
	GetEventsAfterID(id string) []Event
}
