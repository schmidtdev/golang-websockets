package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"schmidtdev/golang-websockets/types"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		//origin := r.Header.Get("Origin")
		//return origin == "http://localhost"
		return true
	},
}

var clients = make(map[*websocket.Conn]bool) // Connected clients
var broadcast = make(chan []byte)            // Broadcast channel
var mutex = &sync.Mutex{}

func WsHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// Upgrade HTTP connection to a WS connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}
	defer conn.Close()

	// Extract IP address
	ip := strings.Split(r.RemoteAddr, ":")[0]

	mutex.Lock()
	clients[conn] = true
	mutex.Unlock()

	go HandleMessages(db, ip)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			mutex.Lock()
			delete(clients, conn)
			mutex.Unlock()
			break
		}

		var receivedMsg types.Message
		err = json.Unmarshal(message, &receivedMsg)
		if err != nil {
			fmt.Println("Error decoding message:", err)
			continue
		}

		broadcast <- message
	}
}

func HandleMessages(db *gorm.DB, ip string) {
	for {
		// Grab next message
		message := <-broadcast

		var receivedMsg types.Message
		err := json.Unmarshal(message, &receivedMsg)
		if err != nil {
			fmt.Println("Error decoding message:", err)
			continue
		}

		var response = types.Message{}

		if receivedMsg.Type == "subscription" {
			var channel types.Channel
			err := db.Where("name = ?", receivedMsg.Channel).First(&channel).Error

			if err != nil {
				db.Create(&types.Channel{Name: receivedMsg.Channel})
			}

			var subscription = types.Subscription{}
			err = db.Where("channel_id = ? AND ip = ?", channel.ID, ip).First(&subscription).Error

			if err != nil {
				db.Create(&types.Subscription{
					ChannelID: channel.ID,
					IP:        ip,
				})

				response.Content = "You have been subscribed to the channel " + receivedMsg.Channel
				response.Type = "subscription"
				response.Channel = receivedMsg.Channel
				response.IP = ip
			} else {
				response.Content = "You are already subscribed to the channel " + receivedMsg.Channel
				response.Type = "subscription"
				response.Channel = receivedMsg.Channel
				response.IP = ip
			}
		} else if receivedMsg.Type == "unsubscription" {
			var channel types.Channel
			err := db.Where("name = ?", receivedMsg.Channel).First(&channel).Error

			if err != nil {
				db.Create(&types.Channel{Name: receivedMsg.Channel})
			}

			var subscription = types.Subscription{}
			err = db.Where("channel_id = ? AND ip = ?", channel.ID, ip).First(&subscription).Error

			if err != nil {
				response.Content = "You have already unsubscribed."
				response.Type = "unsubscription"
				response.Channel = receivedMsg.Channel
				response.IP = ip
			} else {
				err := db.Where("channel_id = ? AND ip = ?", channel.ID, ip).Delete(&subscription).Error

				if err != nil {
					response.Content = "It's not possible to unsubscribe right now..."
					response.Type = "unsubscription"
					response.Channel = receivedMsg.Channel
					response.IP = ip
				} else {
					response.Content = "You have been unsubscribed."
					response.Type = "unsubscription"
					response.Channel = receivedMsg.Channel
					response.IP = ip
				}
			}
		} else {
			response.Content = "You have been unsubscribed from the channel"
			response.Type = "unsubscription"
			response.Channel = "general"
			response.IP = ip
		}

		jsonResponse, err := json.Marshal(response)
		if err != nil {
			fmt.Println("Error marshalling response:", err)
			continue
		}

		// Send to all
		mutex.Lock()
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, jsonResponse)
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()
	}
}
