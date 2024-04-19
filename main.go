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

	}
}

func (s *Server)readLoop(c net.Conn){
	for {
		n, err:= 
	}
}

func main() {

}
