package websockets

import (
	"bytes"
	"context"
	"fmt"
	"stockmarket/database"
	models "stockmarket/models"
	websocketModels "stockmarket/models/websockets"
	gameTempl "stockmarket/templates/games"
)

var hub *websocketModels.Hub

func RunHub() {
	for {
		select {
		case client := <-hub.Register:
			hub.Clients[client] = true
		case client := <-hub.Unregister:
			if _, ok := hub.Clients[client]; ok {
				fmt.Println("client unregistered, deleting and closing", client.Player.User.Name)
				delete(hub.Clients, client)
				close(client.Send)
			}
		case broadcastMessage := <-hub.Broadcast:

			buffer := broadcastMessage.Buffer
			broadcast_game := broadcastMessage.Game
			message := broadcastMessage.Message

			fmt.Println("broadcasted! ", broadcast_game.ID, message)

			for client := range hub.Clients {
				fmt.Println("checking", client.Player.User.Name)
				// only send message to clients in the same game
				if client.Game.ID != broadcast_game.ID {
					fmt.Println("client not in game", broadcast_game.ID, "instead in", client.Game.ID)
					continue
				}

				// buffer == nil when the template broadcast requires user context.
				// these requests reference the DB so should be used sparingly
				if buffer == nil {
					fmt.Println("creating unique buffer for each client and updating client user")

					db := database.GetDb()

					game, err := models.GetGame(client.Game.ID, db)

					if err != nil {
						fmt.Println("could not update game in client")
						continue
					}

					player, err := models.LoadPlayer(client.Player.ID, db)

					if err != nil {
						fmt.Println("could not update player in client")
						continue
					}

					client.Game = game
					client.Player = player

					fmt.Println("client.User updated", client.Player.User.Name)

					if message == "start play" {
						fmt.Println("rendering 'start play' socket for:", client.Player.User.Name, client.Game.ID)

						if err != nil {
							fmt.Println("error getting player from game, perhaps they left and the connection wasn't removed?:", err)
							continue
						}

						boardDisplay := gameTempl.PlayingSocket(client.Game, client.Player)

						buffer = &bytes.Buffer{}
						boardDisplay.Render(context.Background(), buffer)

						// fmt.Println("client: ", client.UserID)
					}

				}

				select {
				case client.Send <- buffer:
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
