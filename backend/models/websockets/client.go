// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websockets

import (
	"bytes"

	models "stockmarket/models"

	gorrilaws "github.com/gorilla/websocket"
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	Hub *Hub

	// The websocket connection.
	Conn *gorrilaws.Conn

	// Buffered channel of outbound messages.
	Send chan *bytes.Buffer

	Game   models.Game
	Player *models.Player
}

type BroadcastMessage struct {
	Game    models.Game
	Buffer  *bytes.Buffer
	Message string
}

func NewClient(hub *Hub, conn *gorrilaws.Conn, player *models.Player, game models.Game) *Client {
	return &Client{
		Hub:    hub,
		Conn:   conn,
		Player: player,
		Game:   game,

		Send: make(chan *bytes.Buffer),
	}
}
