package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"schmidtdev/golang-websockets/types"
	"schmidtdev/golang-websockets/webgr"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		//origin := r.Header.Get("Origin")
		//return origin == "http://integracoes.webgrapp.com.br"
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

		switch receivedMsg.Type {
		case "subscription":
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
		case "unsubscription":
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
		case "get:subscriptions":
			var subscriptions []types.Subscription
			err := db.Where("ip = ?", ip).Joins("Channel").Find(&subscriptions).Error

			if err != nil {
				response.Content = "It's not possible to get your subscriptions right now..."
				response.Type = "response"
				response.Channel = "subscriptions"
				response.IP = ip
			} else {
				response.Type = "response"
				response.Channel = "subscriptions"
				response.IP = ip

				subscriptionList := []string{}
				for _, subscription := range subscriptions {
					subscriptionList = append(subscriptionList, subscription.Channel.Name)
				}
				response.Content = strings.Join(subscriptionList, "\n")
			}
		case "post:webgr_pedidos":
			var receivedMsg types.Message
			err := json.Unmarshal(message, &receivedMsg)
			if err != nil {
				fmt.Println("Error decoding message:", err)
				continue
			}

			var pedidos []webgr.WebgrPedido
			err = json.Unmarshal([]byte(receivedMsg.Content), &pedidos)
			if err != nil {
				fmt.Println("Error decoding pedidos:", err)
				continue
			}

			for _, pedido := range pedidos {
				err = db.Where("cdpedido = ?", pedido.Cdpedido).First(&webgr.WebgrPedido{}).Error
				if err == nil {
					err = db.Model(&pedido).Updates(map[string]interface{}{
						"status":                     pedido.Status,
						"obsinterna":                 pedido.Obsinterna,
						"pedaviso":                   pedido.Pedaviso,
						"cdaut":                      pedido.Cdaut,
						"ftaut":                      pedido.Ftaut,
						"nrnfe":                      pedido.Nrnfe,
						"nrcfe":                      pedido.Nrcfe,
						"erpstatus":                  pedido.Erpstatus,
						"motivocancelamento":         pedido.Motivocancelamento,
						"erpstatusdesc":              pedido.Erpstatusdesc,
						"latitude":                   pedido.Latitude,
						"longitude":                  pedido.Longitude,
						"percentualcomissao":         pedido.Percentualcomissao,
						"vrcomissao":                 pedido.Vrcomissao,
						"vrbruto":                    pedido.Vrbruto,
						"vrdesconto":                 pedido.Vrdesconto,
						"vrdescontototal":            pedido.Vrdescontototal,
						"vrdescontofinal":            pedido.Vrdescontofinal,
						"percentualdescontofinal":    pedido.Percentualdescontofinal,
						"tipodesconto":               pedido.Tipodesconto,
						"vrdescontocondpgto":         pedido.Vrdescontocondpgto,
						"percentualdescontocondpgto": pedido.Percentualdescontocondpgto,
						"tipodescontocondpgto":       pedido.Tipodescontocondpgto,
						"vrliquidototal":             pedido.Vrliquidototal,
						"vrtotal":                    pedido.Vrtotal,
						"cncondicaopgto":             pedido.Cncondicaopgto,
						"cnformapgto":                pedido.Cnformapgto,
					}).Error

					if err != nil {
						fmt.Println("Error updating pedido:", err)
						continue
					} else {
						continue
					}
				} else {
					err = db.Create(&pedido).Error
					if err != nil {
						fmt.Println("Error inserting pedido:", err)
						continue
					}
				}
			}

			response.Content = "Orders have been inserted successfully."
			response.Type = "response"
			response.Channel = "webgr_pedidos"
			response.IP = ip
		case "get:webgr_pedidos":
			dsn := fmt.Sprintf(
				"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
				os.Getenv("WEBGR_HOST"), os.Getenv("WEBGR_USER"), os.Getenv("WEBGR_PASSWORD"), os.Getenv("WEBGR_DB"), os.Getenv("WEBGR_PORT"), os.Getenv("WEBGR_SSL_MODE"),
			)
			db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
				Logger: logger.Default.LogMode(logger.Silent),
			})

			if err != nil {
				response.Content = "It's not possible to connect to the database right now..."
				response.Type = "response"
				response.Channel = "pedidos"
				response.IP = ip
			}

			var pedidos []webgr.WebgrPedido
			err = db.Where("cnempresa = ? AND status = '70' ORDER BY cdpedido DESC LIMIT 40", 71).Find(&pedidos).Error

			if err != nil {
				response.Content = "It's not possible to get the orders right now..."
				response.Type = "response"
				response.Channel = "pedidos"
				response.IP = ip
			} else {
				response.Type = "response"
				response.Channel = "pedidos"
				response.IP = ip

				pedidosList := []map[string]interface{}{}
				for _, pedido := range pedidos {
					pedidosList = append(pedidosList, map[string]interface{}{
						"cdpedido": pedido.Cdpedido,
						"vrtotal":  pedido.Vrtotal,
					})
				}
				content, err := json.Marshal(pedidosList)
				if err != nil {
					response.Content = "Error marshalling pedidos list"
				} else {
					response.Content = string(content)
				}
			}
		default:
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
