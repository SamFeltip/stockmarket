package router

import (
	"log"
	"net/http"
	"stockmarket/middleware"
	"stockmarket/models"
	"stockmarket/websockets"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// serveWs handles websocket requests from the peer.
func serveWs(hub *websockets.Hub, c *gin.Context) {
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

	client := websockets.NewClient(hub, conn, userID, gameID)
	client.Hub.Register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.WritePump()
	go client.ReadPump()
}

func CreateWebsocketRoutes(db *gorm.DB, r *gin.Engine) {

	hub := websockets.NewHub()
	go hub.Run()

	r.GET("/ws",
		func(c *gin.Context) { middleware.RequireAuthWebsocket(c, db) },
		func(c *gin.Context) { serveWs(hub, c) })
}
