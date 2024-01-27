package router

import (
	"log"
	"net/http"
	"stockmarket/websockets"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// serveWs handles websocket requests from the peer.
func serveWs(hub *websockets.Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// page

	client := websockets.NewClient(hub, conn)
	client.Hub.Register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.WritePump()
	go client.ReadPump()
}

func CreateWebsocketRoutes(db *gorm.DB, r *gin.Engine) {

	hub := websockets.NewHub()
	go hub.Run()

	r.GET("/ws", func(c *gin.Context) {
		serveWs(hub, c.Writer, c.Request)
	})
}
