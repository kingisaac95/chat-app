package main

/*
	forward: a channel that holds incoming messages
	that should be forwarded to the other clients.
*/
type room struct {
	forward chan []byte
}
