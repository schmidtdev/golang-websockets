package handlers

import (
	"fmt"

	"github.com/gorilla/websocket"
)

func HandleConnection(conn *websocket.Conn) {
	// Listen
	for {
		// Read message
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}
		fmt.Printf("Received: %s", message)

		// Echo message back
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			fmt.Println("Error writing message:", err)
			break
		}
	}
}
