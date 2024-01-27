package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	page "stockmarket/templates"

	"github.com/a-h/templ"
)

type Request struct {
	Template string `json:"template"`
	Message  string `json:"message"`
	Headers  struct {
		HXRequest     string `json:"HX-Request"`
		HXTrigger     string `json:"HX-Trigger"`
		HXTriggerName string `json:"HX-Trigger-Name"`
		HXTarget      string `json:"HX-Target"`
		HXCurrentURL  string `json:"HX-Current-URL"`
	} `json:"HEADERS"`
}

func CreateWebsocketBuffer(message []byte) (*bytes.Buffer, error) {

	// convery msg to json
	var request Request
	err := json.Unmarshal(message, &request)
	if err != nil {
		return nil, fmt.Errorf("failed to parse message as request: %v", string(message))
	}

	var component templ.Component

	// updating game state (buy, sell, pass)
	// joining game
	// leaving game
	switch request.Template {
	case "chat_message":
		component = page.Card(request.Message)
	}

	if component == nil {
		fmt.Println("No component found for request:", request)
		return nil, fmt.Errorf("no component found for request: %v", request)
	}

	buffer := &bytes.Buffer{}
	component.Render(context.Background(), buffer)
	return buffer, nil
}
