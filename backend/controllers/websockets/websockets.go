package websockets

import (
	"log"
	"net/http"
	"stockmarket/models"
	"stockmarket/websockets"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// serveWs handles websocket requests from the peer.
func ServeWs(c *gin.Context) {

	w := c.Writer
	r := c.Request

	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	cu, _ := c.Get("user")
	cg, _ := c.Get("game")

	if cu == nil {
		log.Println("no user found")
		c.JSON(http.StatusBadRequest, gin.H{"error": "no user found in request context"})
		return
	}

	if cg == nil {
		log.Println("no game found")
		c.JSON(http.StatusBadRequest, gin.H{"error": "no game found in request context"})
		return
	}

	userID := cu.(models.User).ID
	gameID := cg.(models.Game).ID

	client := websockets.NewClient(conn, userID, gameID)

	hub := websockets.GetHub()
	hub.Register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.WritePump()
	go client.ReadPump()
}
