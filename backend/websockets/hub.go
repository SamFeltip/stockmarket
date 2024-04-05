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

var gameHub *websocketModels.Hub

func RunGameHub() {
	for {
		select {
		case client := <-gameHub.Register:
			gameHub.Clients[client] = true
		case client := <-gameHub.Unregister:
			if _, ok := gameHub.Clients[client]; ok {
				fmt.Println("client unregistered, deleting and closing", client.CurrentPlayerID)
				delete(gameHub.Clients, client)
				close(client.Send)
			}
		case broadcastMessage := <-gameHub.Broadcast:

			buffer := broadcastMessage.Buffer
			broadcast_game_id := broadcastMessage.GameID
			message := broadcastMessage.Message

			fmt.Println("broadcasted! ", broadcast_game_id, message)

			db := database.GetDb()

			for client := range gameHub.Clients {
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

						specialPlayerInsights, err := models.LoadSpecialPlayerInsights(current_player.ID, db)

						if err != nil {
							fmt.Println("error loading special insights:", err)
							continue
						}

						latestFeedItem, err := models.LoadLatestFeedItem(game.ID, db)

						if err != nil {
							fmt.Println("error loading latest insight:", err)
							continue
						}

						boardDisplay := gameTempl.PlayingSocket(game, current_player, players, specialPlayerInsights, latestFeedItem)

						buffer = &bytes.Buffer{}
						boardDisplay.Render(context.Background(), buffer)

					} else if message == "special insights" {
						fmt.Println("broadcasting show special insights")

						specialInsights, err := models.LoadSpecialInsights(game.ID, db)

						if err != nil {
							fmt.Println("could not get game insights", err)
							continue
						}

						gameStockDisplays, err := models.LoadGameStockDisplays(game.ID, false, db)

						if err != nil {
							fmt.Println("could not load game stock displays", err)
							continue
						}

						displayGameStocks, err := models.LoadGameStockDisplays(game.ID, true, db)

						if err != nil {
							fmt.Println("could not load game stock displays", err)
							continue
						}

						playerDisplays, err := models.LoadPlayerDisplays(game.ID, db)

						if err != nil {
							fmt.Println("could not load player displays", err)
							continue
						}

						currentPlayerDisplay, err := models.LoadPlayerDisplay(current_player.ID, db)

						if err != nil {
							fmt.Println("could not load current player display", err)
							continue
						}

						specialInsightsDisplay := gameTempl.SpecialInsightsSocket(game.ID, specialInsights, gameStockDisplays, displayGameStocks, playerDisplays, currentPlayerDisplay)

						buffer := &bytes.Buffer{}
						specialInsightsDisplay.Render(context.Background(), buffer)

					}
				}

				select {
				case client.Send <- buffer:
				default:
					close(client.Send)
					delete(gameHub.Clients, client)
				}

				buffer = broadcastMessage.Buffer

			}
		}
	}
}

var playerInsightHub *websocketModels.Hub

func RunPlayerInsightHub() {
	for {
		select {
		case client := <-playerInsightHub.Register:
			playerInsightHub.Clients[client] = true
		case client := <-playerInsightHub.Unregister:
			if _, ok := playerInsightHub.Clients[client]; ok {
				fmt.Println("client unregistered, deleting and closing", client.CurrentPlayerID)
				delete(playerInsightHub.Clients, client)
				close(client.Send)
			}
		case broadcastMessage := <-playerInsightHub.Broadcast:

			buffer := broadcastMessage.Buffer
			broadcast_game_id := broadcastMessage.GameID
			message := broadcastMessage.Message

			fmt.Println("broadcasted! ", broadcast_game_id, message)

			for client := range playerInsightHub.Clients {
				fmt.Println("checking", client.CurrentPlayerID)
				// only send message to clients in the same game
				if client.GameID != broadcast_game_id {
					fmt.Println("client not in game", broadcast_game_id, "instead in", client.GameID)
					continue
				}
				fmt.Println("client in game", broadcast_game_id)

				select {
				case client.Send <- buffer:
				default:
					close(client.Send)
					delete(playerInsightHub.Clients, client)
				}

				buffer = broadcastMessage.Buffer

			}
		}
	}
}

func NewGameHub() *websocketModels.Hub {

	gameHub = &websocketModels.Hub{
		Broadcast:  make(chan *websocketModels.BroadcastMessage),
		Register:   make(chan *websocketModels.Client),
		Unregister: make(chan *websocketModels.Client),
		Clients:    make(map[*websocketModels.Client]bool),
	}
	return gameHub
}

func NewPlayerInsightHub() *websocketModels.Hub {

	playerInsightHub = &websocketModels.Hub{
		Broadcast:  make(chan *websocketModels.BroadcastMessage),
		Register:   make(chan *websocketModels.Client),
		Unregister: make(chan *websocketModels.Client),
		Clients:    make(map[*websocketModels.Client]bool),
	}
	return playerInsightHub
}

func InitializeGameHub() *websocketModels.Hub {
	hub := NewGameHub()

	go RunGameHub()

	return hub
}

func InitializeGameClosedHub() *websocketModels.Hub {
	hub := NewPlayerInsightHub()

	go RunPlayerInsightHub()

	return hub
}

func GetGameHub() *websocketModels.Hub {
	return gameHub
}

func GetPlayerInsightHub() *websocketModels.Hub {
	return playerInsightHub
}
