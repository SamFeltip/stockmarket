package websockets

import (
	"fmt"
	"log"
	"net/http"
	"stockmarket/models"
	websocketModels "stockmarket/models/websockets"
	"stockmarket/websockets"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// serveWs handles websocket requests from the peer.
func ServeWs(c *gin.Context) (int, gin.H) {

	w := c.Writer
	r := c.Request

	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return http.StatusInternalServerError, gin.H{"error": "could not upgrade websocket"}
	}

	cp, player_exists := c.Get("player")

	gameID := c.Param("gameID")

	if !player_exists {
		log.Println("websocket: no user found")
		return http.StatusBadRequest, gin.H{"error": "no user found in request context"}
	}

	player := cp.(models.Player)

	hub := websockets.GetHub()

	if err != nil {
		fmt.Println("error setting active game:", err)
	}

	client := websocketModels.NewClient(hub, conn, player.ID, gameID)

	fmt.Println("registering new client", player.User.Name, gameID)
	hub.Register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go websockets.WritePump(client)
	go websockets.ReadPump(client)

	return http.StatusOK, gin.H{"message": "websocket connection established"}
}
