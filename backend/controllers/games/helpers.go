package games

import (
	"bytes"
	"context"
	"fmt"
	"stockmarket/models"
	websocketModels "stockmarket/models/websockets"
	templates "stockmarket/templates/games"
	"stockmarket/websockets"

	"github.com/a-h/templ"
	"gorm.io/gorm"
)

func BroadcastUpdatePlayersList(gameID string, userCardList templ.Component) error {

	buffer := &bytes.Buffer{}
	userCardList.Render(context.Background(), buffer)

	broadcastMessage := websocketModels.BroadcastMessage{
		GameID: gameID,
		Buffer: buffer,
	}

	hub := websockets.GetHub()
	hub.Broadcast <- &broadcastMessage //send a html template on the hub's broadcast channel

	return nil
}

func BroadcastUpdatePeriodCount(game models.Game) error {

	periodCountDisplay := templates.PeriodCountOptionsSocket(game)

	buffer := &bytes.Buffer{}
	periodCountDisplay.Render(context.Background(), buffer)

	broadcastMessage := websocketModels.BroadcastMessage{
		GameID: game.ID,
		Buffer: buffer,
	}

	hub := websockets.GetHub()
	hub.Broadcast <- &broadcastMessage //send a html template on the hub's broadcast channel

	return nil
}

func BroadcastUpdatePlayBoard(gameID string) error {

	fmt.Println("broadcasting show board: capturing playing socket template")

	broadcastMessage := websocketModels.BroadcastMessage{
		GameID:  gameID,
		Buffer:  nil,
		Message: "game board",
	}

	fmt.Println("broadcasting show board: sending playing socket template")
	hub := websockets.GetHub()
	hub.Broadcast <- &broadcastMessage //send a html template on the hub's broadcast channel
	return nil
}

func BroadcastGameClosed(gameInsights []models.GameInsight, gameID string) error {
	fmt.Println("broadcasting market closed")

	marketClosedDisplay := templates.ClosedSocket(gameInsights, game.GameStocks, game.Players)

	buffer := &bytes.Buffer{}
	marketClosedDisplay.Render(context.Background(), buffer)

	broadcastMessage := websocketModels.BroadcastMessage{
		GameID:  gameID,
		Buffer:  buffer,
		Message: "market closed",
	}

	hub := websockets.GetHub()
	hub.Broadcast <- &broadcastMessage //send a html template on the hub's broadcast channel

	return nil
}

func CheckForMarketClose(gameID string, db *gorm.DB) (templ.Component, error) {

	game, err := models.LoadGameDisplay(gameID, db)

	if err != nil {
		fmt.Println("could not load game display", err)
		return nil, err
	}

	if game.CurrentTurn() <= game.PlayerCount*3 {
		fmt.Println("market not closed")
		return nil, nil
	}

	err = db.Model(models.Game{}).Where("id = ?", gameID).Update("status", string(models.Closed)).Error

	if err != nil {
		fmt.Println("could not close game", err)
		return templates.Error(err), err
	}

	gameInsights, err := models.GetGameInsights(gameID, db)

	if err != nil {
		fmt.Println("could not get game insights", err)
		return templates.Error(err), err
	}

	err = BroadcastGameClosed(gameInsights, gameID)

	if err != nil {
		fmt.Println("could not broadcast market closed", err)
		return templates.Error(err), err
	}

	return templates.Loading(), nil
}
