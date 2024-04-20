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
	peerMap   map[net.Addr]*Client
	mutex     sync.Mutex
	cchan     chan *Client
}

func newServer(listenAdr string) *Server {
	return &Server{
		listenAdr: listenAdr,
		quitch:    make(chan struct{}),
		msgch:     make(chan Message, 10),
		peerMap:   make(map[net.Addr]*Client),
		cchan:     make(chan *Client, 2),
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
	go s.handleMsg()
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

func (s *Server) handleClient(conn net.Conn) {
	conn.Write([]byte("ENTER_NAME"))
	buff := make([]byte, 255)
	n, err := conn.Read(buff)
	if err != nil {
		log.Println("read nickname err:", err)
		conn.Close()
	}
	name := string(buff[:n])
	log.Println(name)
	s.addClient(conn, name)
	s.readLoop(conn)
}

func (s *Server) readLoop(conn net.Conn) {
	buff := make([]byte, 2048)
	for {
		n, err := conn.Read(buff)
		if err != nil {
			log.Println("read err:", err)
			conn.Close()
			continue
		}
		s.msgch <- Message{
			client:  conn,
			payload: buff[:n],
		}

	}
}
func (s *Server) addClient(conn net.Conn, name string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.peerMap[conn.RemoteAddr()] = &Client{
		name: name,
		conn: conn,
	}
}

func (s *Server) broadcastAll(msg string) {
	for _, client := range s.peerMap {
		client.conn.Write([]byte(msg))
	}
}

type Client struct {
	name string
	conn net.Conn
}

type Message struct {
	client  net.Conn
	payload []byte
}

func (s *Server) handleMsg() {
	for msg := range s.msgch {
		log.Println(msg.client, "wrote: ", msg.payload)
		s.broadcastAll(string(msg.payload))
	}
}

func main() {
	server := newServer(":3000")
	go server.handleMsg()
	log.Println("Listening...")
	log.Fatal((server.Run()))
}
