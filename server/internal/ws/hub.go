package ws

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"github.com/puremike/online_auction_api/internal/models"
)

type Hub struct {
	// connected clients
	// map[string]map[*websocket.Conn]bool where outer key is userID, inner map is connections for that user
	Clients    map[string]map[*websocket.Conn]bool
	Broadcast  chan []byte  // handle incoming messages from clients
	Register   chan *Client // register requests from clients
	Unregister chan *Client // unregister requests from clients

	AuctionUpdates      chan *models.AuctionUpdateEvent
	NotificationUpdates chan *models.NotificationEvent
}

func NewHub() *Hub {
	return &Hub{
		Clients:             make(map[string]map[*websocket.Conn]bool),
		Broadcast:           make(chan []byte),
		Register:            make(chan *Client),
		Unregister:          make(chan *Client),
		AuctionUpdates:      make(chan *models.AuctionUpdateEvent),
		NotificationUpdates: make(chan *models.NotificationEvent),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			if _, ok := h.Clients[client.userID]; !ok {
				h.Clients[client.userID] = make(map[*websocket.Conn]bool)
			}
			h.Clients[client.userID][client.conn] = true
			log.Printf("Client %s connected (Conn:%p)", client.userID, client.conn)

		case client := <-h.Unregister:
			if connections, ok := h.Clients[client.userID]; ok {
				if _, ok := connections[client.conn]; ok {
					delete(connections, client.conn)
					close(client.send) //close send channel to stop writePump
					log.Printf("Client unregisterd: %s (Conn: %p)", client.userID, client.conn)
					if len(connections) == 0 {
						delete(h.Clients, client.userID) // remove user if no active connetions
					}
				}
			}

		case message := <-h.Broadcast:
			for _, userConnections := range h.Clients {
				for conn := range userConnections {
					if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
						log.Printf("Failed to send message to client: %v (Conn: %p)", err, conn)
					}
				}
			}

		case auctionUpdate := <-h.AuctionUpdates:
			log.Printf("Received auction update for AuctionID: %s, CurrentPrice: %.2f", auctionUpdate.ID, auctionUpdate.CurrentPrice)

			// marshal the event to JSON
			jsonData, err := json.Marshal(auctionUpdate)
			if err != nil {
				log.Printf("Failed to marshal auction update event: %v", err)
				continue
			}

			// broadcast the event to all connected clients
			for _, userConnections := range h.Clients {
				for conn := range userConnections {
					if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
						log.Printf("Failed to send auction update to client: %v (Conn: %p)", err, conn)
					}
				}
			}

		case notificationUpdates := <-h.NotificationUpdates:
			log.Printf("Received notification update for UserID: %s, Message: %s", notificationUpdates.UserID, notificationUpdates.Message)

			jsonData, err := json.Marshal(notificationUpdates)
			if err != nil {
				log.Printf("Failed to marshal notification update event: %v", err)
				continue
			}

			if connections, ok := h.Clients[notificationUpdates.UserID]; ok {
				for conn := range connections {
					if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
						log.Printf("Failed to send notification update to client: %v (Conn: %p)", err, conn)
					}
				}
			}
		}
	}
}
