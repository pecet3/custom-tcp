package main

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	listenAdr string
	ln        net.Listener
	quitch    chan struct{}
	msgch     chan Message
}

type Message struct {
	client  net.Conn
	payload []byte
}

func NewServer(listenAdr string) *Server {
	s := &Server{
		listenAdr: listenAdr,
		quitch:    make(chan struct{}),
		msgch:     make(chan Message, 10),
	}
	return s
}

func (s *Server) Run() error {
	ln, err := net.Listen("tcp", s.listenAdr)
	if err != nil {
		return err
	}
	defer ln.Close()
	s.ln = ln

	go s.acceptLoop()

	<-s.quitch
	defer close(s.msgch)
	return nil
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			log.Println("acceptLoop err: ", err)
			continue
		}
		log.Println(conn)
		go s.readLoop(conn)

	}
}

func (s *Server) readLoop(conn net.Conn) {
	buff := make([]byte, 2048)
	for {
		n, err := conn.Read(buff)
		if err != nil {
			log.Println("read err:", err)
			continue
		}
		payload := buff[0:n]
		log.Println(n)
		s.msgch <- Message{
			client:  conn,
			payload: payload,
		}
	}

}

func (s *Server) handleMsg() {
	for msg := range s.msgch {
		fmt.Print(msg.client.LocalAddr())
		fmt.Println(string(msg.payload))
	}
}

func main() {
	server := NewServer(":3000")

	go server.handleMsg()
	log.Println("Listening...")
	log.Fatal((server.Run()))
}
