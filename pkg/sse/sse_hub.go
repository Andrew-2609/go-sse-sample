package sse

type SSEHub struct {
	clients    map[chan Event]struct{}
	Register   chan chan Event
	Unregister chan chan Event
	Broadcast  chan Event
}

func NewSSEHub() *SSEHub {
	h := &SSEHub{
		clients:    make(map[chan Event]struct{}),
		Register:   make(chan chan Event),
		Unregister: make(chan chan Event),
		Broadcast:  make(chan Event),
	}

	go h.run()

	return h
}

func (h *SSEHub) run() {
	for {
		select {
		case c := <-h.Register:
			h.clients[c] = struct{}{}
		case c := <-h.Unregister:
			delete(h.clients, c)
			close(c)
		case event := <-h.Broadcast:
			for c := range h.clients {
				select {
				case c <- event:
				default:
					// slow client -> drop it
					delete(h.clients, c)
					close(c)
				}
			}
		}
	}
}
