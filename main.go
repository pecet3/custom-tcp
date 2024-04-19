package main

import (
	"log"
	"net"
)

type Server struct {
	listenAdr string
	ln        net.Listener
	quitch    chan struct{}
}

func NewServer(listenAdr string) *Server {
	s := &Server{
		listenAdr: listenAdr,
		quitch:    make(chan struct{}),
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
		log.Println("n = ", n)

		msg := buff[:n]
		log.Println(string(msg))
	}

}

func main() {
	server := NewServer(":3000")
	server.Run()
}
