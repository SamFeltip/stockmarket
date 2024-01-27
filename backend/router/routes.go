package router

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"stockmarket/models"
	page "stockmarket/templates"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wshandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("connection made")
	conn, err := wsupgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println("Failed to set websocket upgrade:", err)
		return
	}

	for {
		_, msg, err := conn.ReadMessage()
		fmt.Println("socket: ", string(msg))
		if err != nil {
			fmt.Printf("Failed to read message:", err)
			break
		}
		// c.HTML(http.StatusOK, "<div>something</div>", nil)

		// convery msg to json
		var response Response
		err = json.Unmarshal(msg, &response)
		if err != nil {
			fmt.Println("Failed to parse message:", err)
			continue
		}

		cardComponent := page.Card(response.ChatMessage)
		buffer := &bytes.Buffer{}
		cardComponent.Render(context.Background(), buffer)

		// websocket.TextMessage.Send(websocket.TextMessage, buffer.Bytes())
		err = conn.WriteMessage(websocket.TextMessage, buffer.Bytes())
		if err != nil {
			fmt.Println("Failed to send message:", err)
			continue
		}

	}

	fmt.Println("connection completed")
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{hub: hub, conn: conn, send: make(chan *bytes.Buffer)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}

func SetupRoutes(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.LoadHTMLFiles("sockets.html")

	hub := newHub()
	go hub.run()

	r.GET("/ws", func(c *gin.Context) {
		// wshandler(c.Writer, c.Request)
		serveWs(hub, c.Writer, c.Request)
	})

	//router.LoadHTMLFiles("templates/template1.html", "templates/template2.html")
	r.GET("/sockets", func(c *gin.Context) {
		c.HTML(http.StatusOK, "sockets.html", nil)
	})

	CreateAuthRoutes(db, r)

	CreatePageRoutes(db, r)
	CreateUserRoutes(db, r)
	CreateGameRoutes(db, r)

	return r
}

func RenderWithTemplate(pageComponent templ.Component, title string, c *gin.Context) {

	cu, _ := c.Get("user")

	if cu == nil {
		cu = models.User{}
	}

	user := cu.(models.User)

	ctx := context.WithValue(context.Background(), page.CurrentUser, user)

	baseComponent := page.Base(title, pageComponent, c)
	baseComponent.Render(ctx, c.Writer)
}
