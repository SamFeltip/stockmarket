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
				fmt.Println("client unregistered, deleting and closing", client.CurrentPlayerID)
				delete(hub.Clients, client)
				close(client.Send)
			}
		case broadcastMessage := <-hub.Broadcast:

			buffer := broadcastMessage.Buffer
			broadcast_game_id := broadcastMessage.GameID
			message := broadcastMessage.Message

			fmt.Println("broadcasted! ", broadcast_game_id, message)

			db := database.GetDb()

			for client := range hub.Clients {
				fmt.Println("checking", client.CurrentPlayerID)
				// only send message to clients in the same game
				if client.GameID != broadcast_game_id {
					fmt.Println("client not in game", broadcast_game_id, "instead in", client.GameID)
					continue
				}
				fmt.Println("client in game", broadcast_game_id)

				// buffer == nil when the template broadcast requires user context.
				// these requests reference the DB so should be used sparingly
				if buffer == nil {
					fmt.Println("creating unique buffer for each client and updating client user")

					game, err := models.LoadGameDisplay(client.GameID, db)

					if err != nil {
						fmt.Println("could not update game in client")
						continue
					}

					current_player, err := models.LoadCurrentPlayerDisplay(client.CurrentPlayerID, db)

					if err != nil {
						fmt.Println("could not update player in client")
						continue
					}

					fmt.Println("client.User updated", client.CurrentPlayerID)

					if message == "game board" {
						fmt.Println("rendering game board socket for:", client.CurrentPlayerID, client.GameID)

						if err != nil {
							fmt.Println("error getting player from game, perhaps they left and the connection wasn't removed?:", err)
							continue
						}

						players, err := models.LoadPlayerDisplays(game.ID, db)

						if err != nil {
							fmt.Println("error getting players from game, perhaps they left and the connection wasn't removed?:", err)
							continue
						}

						boardDisplay := gameTempl.PlayingSocket(game, current_player, players)

						buffer = &bytes.Buffer{}
						boardDisplay.Render(context.Background(), buffer)

					}
				}

				select {
				case client.Send <- buffer:
				default:
					close(client.Send)
					delete(hub.Clients, client)
				}

				buffer = broadcastMessage.Buffer

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
