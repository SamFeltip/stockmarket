package games

import (
	"bytes"
	"context"
	"fmt"
	"stockmarket/models"
	websocketModels "stockmarket/models/websockets"
	templates "stockmarket/templates/games"
	userTemplates "stockmarket/templates/users"
	"stockmarket/websockets"
	"strconv"

	"github.com/a-h/templ"
	"gorm.io/gorm"
)

func BroadcastUpdatePlayersList(game *models.Game) error {

	userCardList := userTemplates.CardListSocket(game.Players)

	buffer := &bytes.Buffer{}
	userCardList.Render(context.Background(), buffer)

	if len(game.Players) == 0 {
		fmt.Println("no players to broadcast")
		return nil
	}

	broadcastMessage := websocketModels.BroadcastMessage{
		Game:   *game,
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
		Game:   game,
		Buffer: buffer,
	}

	hub := websockets.GetHub()
	hub.Broadcast <- &broadcastMessage //send a html template on the hub's broadcast channel

	return nil
}

func BroadcastUpdatePlayBoard(game models.Game) error {

	fmt.Println("broadcasting show board: capturing playing socket template")

	broadcastMessage := websocketModels.BroadcastMessage{
		Game:    game,
		Buffer:  nil,
		Message: "game board",
	}

	fmt.Println("broadcasting show board: sending playing socket template")
	hub := websockets.GetHub()
	hub.Broadcast <- &broadcastMessage //send a html template on the hub's broadcast channel
	return nil
}

func BroadcastGameClosed(gameInsights []models.GameInsight, game models.Game) error {
	fmt.Println("broadcasting market closed")

	marketClosedDisplay := templates.ClosedSocket(gameInsights, game.GameStocks, game.Players)

	buffer := &bytes.Buffer{}
	marketClosedDisplay.Render(context.Background(), buffer)

	broadcastMessage := websocketModels.BroadcastMessage{
		Game:    game,
		Buffer:  buffer,
		Message: "market closed",
	}

	hub := websockets.GetHub()
	hub.Broadcast <- &broadcastMessage //send a html template on the hub's broadcast channel

	return nil
}

func CheckForMarketClose(game models.Game, db *gorm.DB) (templ.Component, error) {

	fmt.Println("checking for market close")
	fmt.Println("current turn: " + strconv.Itoa(game.CurrentTurn()))
	fmt.Println("players: " + strconv.Itoa(len(game.Players)*3))

	if game.CurrentTurn() <= len(game.Players)*3 {
		fmt.Println("market not closed")
		return nil, nil
	}

	game.Status = string(models.Closed)
	err := db.Save(&game).Error

	if err != nil {
		fmt.Println("could not close game", err)
		return templates.Error(err), err
	}

	gameInsights, err := models.GetGameInsights(game.ID, db)

	if err != nil {
		fmt.Println("could not get game insights", err)
		return templates.Error(err), err
	}

	err = BroadcastGameClosed(gameInsights, game)

	if err != nil {
		fmt.Println("could not broadcast market closed", err)
		return templates.Error(err), err
	}

	return templates.Loading(), nil
}
