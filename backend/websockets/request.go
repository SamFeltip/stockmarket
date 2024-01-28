package websockets

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	websocketModels "stockmarket/models/websockets"

	"github.com/a-h/templ"
)

func HandleWebsocketRequest(message []byte) (*bytes.Buffer, error) {

	// convery msg to json
	var request websocketModels.HTMXRequest
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
		// component = templates.Card(request.Message)
		fmt.Println("chat message:", request.Message)
	}

	if component == nil {
		fmt.Println("No component found for request:", request)
		return nil, fmt.Errorf("no component found for request: %v", request)
	}

	buffer := &bytes.Buffer{}
	component.Render(context.Background(), buffer)
	return buffer, nil
}
