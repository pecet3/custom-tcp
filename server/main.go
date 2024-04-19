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
	peerMap   map[net.Addr]net.Conn
}

type Message struct {
	client  net.Conn
	payload []byte
}

func NewServer(listenAdr string) *Server {
	return &Server{
		listenAdr: listenAdr,
		quitch:    make(chan struct{}),
		msgch:     make(chan Message, 10),
		peerMap:   make(map[net.Addr]net.Conn),
	}
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
	close(s.msgch)

	return nil
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			log.Println("acceptLoop err: ", err)
			continue
		}
		log.Println("client connected:", conn.RemoteAddr().String())
		s.peerMap[conn.RemoteAddr()] = conn
		log.Println("clients on the server:")
		for client := range s.peerMap {
			fmt.Println("      > ", client.String())
		}
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
		s.msgch <- Message{
			client:  conn,
			payload: buff[:n],
		}

		conn.Write([]byte("thank you for msg\n"))
	}

}

func (s *Server) handleMsg() {
	for msg := range s.msgch {
		log.Println(" // " + msg.client.RemoteAddr().String() + " => " + string(msg.payload))

		for _, conn := range s.peerMap {
			conn.Write([]byte(string(msg.payload)))
		}
	}
}

func main() {
	server := NewServer(":3000")

	go server.handleMsg()
	log.Println("Listening...")
	log.Fatal((server.Run()))
}
