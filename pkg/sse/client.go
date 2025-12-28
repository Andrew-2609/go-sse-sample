package sse

import "time"

type sseClient struct {
	ch             chan Event
	connectedAt    time.Time
	disconnectChan chan struct{}
}

func NewSSEClient(ch chan Event, connectedAt time.Time) *sseClient {
	return &sseClient{
		ch:             ch,
		connectedAt:    connectedAt,
		disconnectChan: make(chan struct{}, 1),
	}
}

func (c *sseClient) CH() chan Event {
	return c.ch
}

func (c *sseClient) Disconnect() <-chan struct{} {
	return c.disconnectChan
}
