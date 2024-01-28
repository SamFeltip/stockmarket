package websockets

import (
	websocketModels "stockmarket/models/websockets"
)

var hub *websocketModels.Hub

func RunHub() {
	for {
		select {
		case client := <-hub.Register:
			hub.Clients[client] = true
		case client := <-hub.Unregister:
			if _, ok := hub.Clients[client]; ok {
				delete(hub.Clients, client)
				close(client.Send)
			}
		case broadcastMessage := <-hub.Broadcast:
			message := broadcastMessage.Buffer
			gameID := broadcastMessage.GameID

			for client := range hub.Clients {

				// only send message to clients in the same game
				if client.GameID != gameID {
					continue
				}

				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(hub.Clients, client)
				}
			}
		}
	}
}

func NewHub() *websocketModels.Hub {

	hub = &websocketModels.Hub{
		Broadcast:  make(chan *websocketModels.BroadcastMessage),
		Register:   make(chan *websocketModels.Client),
		Unregister: make(chan *websocketModels.Client),
		Clients:    make(map[*websocketModels.Client]bool),
	}
	return hub
}

func InitializeHub() *websocketModels.Hub {
	hub := NewHub()

	go RunHub()

	return hub
}

func GetHub() *websocketModels.Hub {
	return hub
}
