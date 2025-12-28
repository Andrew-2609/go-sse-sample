package sse

import "sync"

var (
	sseHubSingleton *SSEHub
	sseHubOnce      sync.Once
)

type SSEHub struct {
	eventStore EventStore
	clients    map[*sseClient]struct{}
	order      []*sseClient
	Register   chan *sseClient
	Unregister chan *sseClient
	Broadcast  chan Event
	maxClients int
}

func InitializeSSEHub(eventStore EventStore, maxClients int) {
	if sseHubSingleton != nil {
		return
	}

	sseHubOnce.Do(func() {
		sseHubSingleton = &SSEHub{
			eventStore: eventStore,
			clients:    make(map[*sseClient]struct{}),
			Register:   make(chan *sseClient),
			Unregister: make(chan *sseClient),
			Broadcast:  make(chan Event),
			maxClients: maxClients,
		}

		go sseHubSingleton.run()
	})
}

func GetSSEHub() *SSEHub {
	if sseHubSingleton == nil {
		panic("SSEHub not initialized")
	}

	return sseHubSingleton
}

func (h *SSEHub) run() {
	for {
		select {
		case c := <-h.Register:
			if len(h.clients) >= h.maxClients {
				oldestClient := h.order[0]
				h.order = h.order[1:]
				delete(h.clients, oldestClient)
				close(oldestClient.ch)
				oldestClient.disconnectChan <- struct{}{}
			}
			h.clients[c] = struct{}{}
			h.order = append(h.order, c)
		case c := <-h.Unregister:
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.ch)
				c.disconnectChan <- struct{}{}
				for i, v := range h.order {
					if v == c {
						h.order = append(h.order[:i], h.order[i+1:]...)
						break
					}
				}
			}
		case event := <-h.Broadcast:
			h.eventStore.StoreEvent(event)
			for c := range h.clients {
				select {
				case c.ch <- event:
				default:
					// slow client -> drop it
					delete(h.clients, c)
					close(c.ch)
					c.disconnectChan <- struct{}{}
				}
			}
		}
	}
}

func (h *SSEHub) GetEventsAfterID(id string) []Event {
	return h.eventStore.GetEventsAfterID(id)
}
