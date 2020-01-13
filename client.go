package main

import (
	"github.com/gorilla/websocket"
)

/*
	forward: a channel that holds incoming messages
	that should be forwarded to the other clients.
*/
type room struct {
	forward chan []byte
}

/*
	client: a single chatting user
	socket: the web socket for this client
	room: the room this client is chatting in
*/
type client struct {
	socket *websocket.Conn
	send   chan []byte
	room   *room
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
