package main

import (
	"fmt"
	"net/http"
	"schmidtdev/golang-websockets/handlers"
)

func main() {
	http.HandleFunc("/ws", handlers.WsHandler)
  go handlers.HandleMessages()

	fmt.Println("Listening on :3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		fmt.Println("Server error:", err)
	}
}
