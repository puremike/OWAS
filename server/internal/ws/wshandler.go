package ws

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/puremike/online_auction_api/contexts"
)

type WSHandler struct {
	hub *Hub
}

func NewWSHandler(hub *Hub) *WSHandler {
	return &WSHandler{
		hub: hub,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // for development only

		// // for production
		// origin := r.Header.Get("Origin")
		// setOrigin := "http://localhost:3000"
		// return origin == setOrigin
	},

	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (ws *WSHandler) ServeWs(c *gin.Context) {

	authUser, err := contexts.GetUserFromContext(c)
	if err != nil {
		log.Printf("websocket upgrade failed: user not authenticated: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("websocket upgrade failed: %v", err)
		return
	}

	client := &Client{
		hub:    ws.hub,
		conn:   conn,
		send:   make(chan []byte, 256),
		userID: authUser.ID,
	}

	client.hub.Register <- client

	go client.readPump()
	go client.writePump()
}
