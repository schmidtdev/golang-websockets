package handlers

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
    //origin := r.Header.Get("Origin")
    //return origin == "http://localhost"
    return true
	},
}

var clients = make(map[*websocket.Conn]bool) // Connected clients
var broadcast = make(chan []byte) // Broadcast channel
var mutex = &sync.Mutex{}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade HTTP connection to a WS connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}
	defer conn.Close()

  mutex.Lock()
  clients[conn] = true
  mutex.Unlock()

  for {
    _, message, err := conn.ReadMessage()
    if err != nil {
      mutex.Lock()
      delete(clients, conn)
      mutex.Unlock()
      break
    }
    broadcast <- message
  }
}

func HandleMessages() {
  for {
    // Grab next message
    message := <-broadcast

    // Send to all
    mutex.Lock()
    for client := range clients {
      err := client.WriteMessage(websocket.TextMessage, message)
      if err != nil {
        client.Close()
        delete(clients, client)
      }
    }
    mutex.Unlock()
  }
}
