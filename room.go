package main

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
}

// listen for messages in any of the three channels
// and execute matching case via select
func (r *room) run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.send)
		case msg := <-r.forward:
			for client := range r.clients {
				client.send <- msg
			}
		}
	}
}
