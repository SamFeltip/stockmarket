// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websockets

import (
	"bytes"

	gorrilaws "github.com/gorilla/websocket"
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	Hub *Hub

	// The websocket connection.
	Conn *gorrilaws.Conn

	// Buffered channel of outbound messages.
	Send chan *bytes.Buffer

	GameID string
	UserID uint
}

type BroadcastMessage struct {
	GameID string
	Buffer *bytes.Buffer
}

func NewClient(hub *Hub, conn *gorrilaws.Conn, userID uint, gameID string) *Client {
	return &Client{
		Hub:    hub,
		Conn:   conn,
		UserID: userID,
		GameID: gameID,

		Send: make(chan *bytes.Buffer),
	}
}
