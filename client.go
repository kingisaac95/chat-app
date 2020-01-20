package main

import (
	"github.com/gorilla/websocket"
)

// client: a single chatting user
type client struct {
	// socket: the web socket for this client
	socket *websocket.Conn
	// send: channel on which the messages are sent
	send chan []byte
	// room: the room this client is chatting in
	room *room
}

func (c *client) read() {
	defer c.socket.Close()
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		c.room.forward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
