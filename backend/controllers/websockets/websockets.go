package websockets

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"stockmarket/database"
	"stockmarket/models"
	websocketModels "stockmarket/models/websockets"
	gameTemplates "stockmarket/templates/games"
	userTemplates "stockmarket/templates/users"
	"stockmarket/websockets"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// serveWs handles websocket requests from the peer.
func LoadPlayers(c *gin.Context) (int, gin.H) {

	w := c.Writer
	r := c.Request

	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, gin.H{"error": "could not upgrade websocket"}
	}

	cu, _ := c.Get("user")
	cg, _ := c.Get("game")

	if cu == nil {
		log.Println("no user found")
		return http.StatusBadRequest, gin.H{"error": "no user found in request context"}
	}

	if cg == nil {
		log.Println("no game found")
		return http.StatusBadRequest, gin.H{"error": "no game found in request context"}
	}

	userID := cu.(models.User).ID
	gameID := cg.(models.Game).ID

	hub := websockets.GetHub()

	client := websocketModels.NewClient(hub, conn, userID, gameID)

	hub.Register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go websockets.WritePump(client)
	go websockets.ReadPump(client)

	return http.StatusOK, gin.H{"message": "websocket connection established"}
}

func UpdateDifficulty(c *gin.Context) (int, gin.H) {
	fmt.Println("updating difficulty")
	cg, _ := c.Get("game")

	if cg == nil {
		log.Println("no game found")
		return http.StatusBadRequest, gin.H{"error": "no game found in request context"}
	}

	gameID := cg.(models.Game).ID
	fmt.Println("game id molded")
	// get difficulty from c
	difficultyStr := c.PostForm("difficulty")

	fmt.Println("getting DB...")
	// update difficulty in db
	db := database.GetDb()
	game, err := models.GetGame(gameID, db)
	if err != nil {
		log.Println("error fetching game:", err)
		return http.StatusInternalServerError, gin.H{"error": "could not fetch game"}
	}

	difficulty, err := strconv.Atoi(difficultyStr)
	if err != nil {
		log.Println("error converting difficulty:", err)
		return http.StatusInternalServerError, gin.H{"error": "could not convert difficulty"}
	}

	game.Difficulty = difficulty
	err = game.UpdateORM(db)
	if err != nil {
		log.Println("error updating game:", err)
		return http.StatusInternalServerError, gin.H{"error": "could not update game"}
	}

	// create DifficultyOptionsSocket template and broadcast it

	difficultyOptions := gameTemplates.DifficultyOptionsSocket(game)

	buffer := &bytes.Buffer{}
	difficultyOptions.Render(context.Background(), buffer)

	broadcastMessage := websocketModels.BroadcastMessage{
		GameID: gameID,
		Buffer: buffer,
	}

	hub := websockets.GetHub()
	hub.Broadcast <- &broadcastMessage //send a html template on the hub's broadcast channel

	return http.StatusOK, gin.H{"message": "difficulty updated"}
}

func BroadcastUpdatePlayersList(players []models.Player) error {

	userCardList := userTemplates.CardListSocket(players)

	buffer := &bytes.Buffer{}
	userCardList.Render(context.Background(), buffer)

	latestPlayer := players[len(players)-1]

	broadcastMessage := websocketModels.BroadcastMessage{
		GameID: latestPlayer.GameID,
		Buffer: buffer,
	}

	hub := websockets.GetHub()
	hub.Broadcast <- &broadcastMessage //send a html template on the hub's broadcast channel

	return nil

}
