package chat

import (
	"chat_app/trace"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type room struct {
	// forward: a channel that holds incoming messages
	// that should be forwarded to the other clients.
	forward chan []byte
	// join: a channel for clients wishing to join the room
	join chan *client
	// leave: a channel for clients wishing to leave the room
	leave chan *client
	// clients: holds all current clients in this room
	clients map[*client]bool
	// Tracer: log all room activities
	Tracer trace.Tracer
}

// NewRoom helper for creating rooms
func NewRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

// listen for messages in any of the three channels
// and execute matching case via select
func (r *room) Run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
			r.Tracer.Trace("New client joined.")
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
			r.Tracer.Trace("Client left.")
		case msg := <-r.forward:
			r.Tracer.Trace("Message received: ", string(msg))
			for client := range r.clients {
				client.send <- msg
				r.Tracer.Trace(" -- sent to client.")
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

// websocket.Upgrader type upgrades our HTTP connection to use web sockets
var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// use Upgrade method from the web socket type to get the socket from a HTTP request
	socket, err := upgrader.Upgrade(w, req, nil)

	if err != nil {
		log.Fatal("ServeHTTP:", err)
		return
	}
	// create a client from the socket
	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}
	// add the client to the current room
	r.join <- client
	defer func() { r.leave <- client }()
	// run the client.write method in a different goroutine(thread)
	go client.write()
	// read from the main thread keeping the connection alive until closed
	client.read()
}
