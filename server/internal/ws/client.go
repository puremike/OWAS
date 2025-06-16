package ws

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	userID string
}

var (
	writeDeadline = time.Second * 10
	readDeadline  = time.Minute * 5
)

// writePump pumps messages from the hub to the websocket connection
// it runs in a goroutine for each client
func (c *Client) writePump() {
	defer func() {
		c.hub.Unregister <- c
	}()

	for message := range c.send {
		if err := c.conn.SetWriteDeadline(time.Now().Add(writeDeadline)); err != nil {
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
		}

		w, err := c.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}

		w.Write(message)

		// flush queued messages
		for i := 0; i < len(c.send); i++ {
			w.Write([]byte("\n"))
			w.Write(message)
		}

		if err := w.Close(); err != nil {
			return
		}
	}
}

func (c *Client) readPump() {
	defer func() {
		c.hub.Unregister <- c
	}()

	c.conn.SetReadLimit(512)
	c.conn.SetReadDeadline(time.Now().Add(readDeadline))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(readDeadline))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error reading message: %v", err)
			}
			break
		}

		log.Printf("Received message from client, %s: %s", c.userID, string(message))

		c.hub.Broadcast <- message
	}
}
