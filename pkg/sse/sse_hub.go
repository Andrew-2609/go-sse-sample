package sse

type SSEHub struct {
	clients    map[*sseClient]struct{}
	order      []*sseClient
	Register   chan *sseClient
	Unregister chan *sseClient
	Broadcast  chan Event
	maxClients int
}

func NewSSEHub(maxClients int) *SSEHub {
	h := &SSEHub{
		clients:    make(map[*sseClient]struct{}),
		Register:   make(chan *sseClient),
		Unregister: make(chan *sseClient),
		Broadcast:  make(chan Event),
		maxClients: maxClients,
	}

	go h.run()

	return h
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
				oldestClient.isDisconnected <- true
			}
			h.clients[c] = struct{}{}
			h.order = append(h.order, c)
		case c := <-h.Unregister:
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.ch)
				c.isDisconnected <- true
				for i, v := range h.order {
					if v == c {
						h.order = append(h.order[:i], h.order[i+1:]...)
						break
					}
				}
			}
		case event := <-h.Broadcast:
			for c := range h.clients {
				select {
				case c.ch <- event:
				default:
					// slow client -> drop it
					delete(h.clients, c)
					close(c.ch)
					c.isDisconnected <- true
				}
			}
		}
	}
}

func (h *SSEHub) Len() int {
	return len(h.clients)
}
