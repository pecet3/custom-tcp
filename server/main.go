package main

import (
	"log"
	"net"
)

type Message struct {
	client  net.Conn
	payload []byte
}

func main() {
	server := NewServer(":3000")

	go server.handleNewClients()

	log.Println("Listening...")
	log.Fatal((server.Run()))
}
