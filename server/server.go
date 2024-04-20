package main

import (
	"log"
	"net"
	"sync"
)

type Server struct {
	listenAdr string
	ln        net.Listener
	quitch    chan struct{}
	msgch     chan Message
	peerMap   map[string]*Client
	mutex     sync.Mutex
	cchan     chan *Client
}

func NewServer(listenAdr string) *Server {
	return &Server{
		listenAdr: listenAdr,
		quitch:    make(chan struct{}),
		msgch:     make(chan Message, 10),
		peerMap:   make(map[string]*Client),
		cchan:     make(chan *Client),
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
			conn.Close()
			continue
		}
		log.Println("client connected:", conn.RemoteAddr().String())

		go s.handleClient(conn)

	}
}

func (s *Server) handleNewClients() {
	for c := range s.cchan {
		s.addClient(c.conn, c.name)
		c.conn.Write([]byte(c.name))
	}
}
