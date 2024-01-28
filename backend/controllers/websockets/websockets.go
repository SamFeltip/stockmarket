package websockets

import (
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
	go client.WritePump()
	go client.ReadPump()

	return http.StatusOK, gin.H{"message": "websocket connection established"}
}
