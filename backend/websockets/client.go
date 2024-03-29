package websockets

import (
	"encoding/json"
	"fmt"
	"stockmarket/database"
	models "stockmarket/models"
	websocketModels "stockmarket/models/websockets"
	"time"

	gorrilaws "github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = gorrilaws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// ReadPump pumps messages from the websocket connection to the Hub.
//
// The application runs ReadPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func ReadPump(c *websocketModels.Client) {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	db := database.GetDb()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			// if message is websocket: close 1001 (going away)
			if !gorrilaws.IsCloseError(err, gorrilaws.CloseGoingAway) {
				fmt.Println("Failed to read message:", err)
				break
			}

			fmt.Println("Websocket closed:", err)

			fmt.Println("details:", c.CurrentPlayerID, c.GameID)
			err := models.PlayerLeft(c.CurrentPlayerID, db)

			if err != nil {
				fmt.Println("error setting player to inactive:", err)
			}

			break
		}

		// convery msg to json
		var request websocketModels.HTMXRequest
		err = json.Unmarshal(message, &request)
		if err != nil {
			fmt.Errorf("failed to parse message as request: %v", string(message))
			break
		}

		fmt.Println("websocket request:", request)
	}
}

// WritePump pumps messages from the hub to the websocket connection.
//
// A goroutine running WritePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func WritePump(c *websocketModels.Client) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case messageBuffer, ok := <-c.Send:
			fmt.Println("starting message")
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				// The hub closed the channel.
				fmt.Println("hub closed the channel")
				c.Conn.WriteMessage(gorrilaws.CloseMessage, []byte{})
				return
			}
			fmt.Println("sending message")
			err := c.Conn.WriteMessage(gorrilaws.TextMessage, messageBuffer.Bytes())

			if err != nil {
				fmt.Println("error writing message:", err)
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(gorrilaws.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
