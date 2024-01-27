package websockets

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan *BroadcastMessage

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan *BroadcastMessage),
		Register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case broadcastMessage := <-h.broadcast:
			message := broadcastMessage.Buffer
			gameID := broadcastMessage.gameID

			for client := range h.clients {

				// only send message to clients in the same game
				if client.gameID != gameID {
					continue
				}

				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
