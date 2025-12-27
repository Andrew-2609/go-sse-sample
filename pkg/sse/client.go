package sse

import "time"

type sseClient struct {
	ch             chan Event
	connectedAt    time.Time
	isDisconnected chan bool
}

func NewSSEClient(ch chan Event, connectedAt time.Time) *sseClient {
	return &sseClient{
		ch:             ch,
		connectedAt:    connectedAt,
		isDisconnected: make(chan bool, 1),
	}
}

func (c *sseClient) CH() chan Event {
	return c.ch
}

func (c *sseClient) IsDisconnected() <-chan bool {
	return c.isDisconnected
}
