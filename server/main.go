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
		go s.handleNewConn(conn)

	}
}

func (s *Server) handleNewConn(conn net.Conn) {
	conn.Write([]byte("ENTER_NAME" + conn.RemoteAddr().String()))
	buff := make([]byte, 255)
	n, err := conn.Read(buff)
	if err != nil {
		log.Println("read nickname err:", err)
		conn.Close()
		return
	}
	name := string(buff[:n])
	c := s.addNewClient(conn, name)
	s.broadcastAll(name + " joins the server")
	c.readLoop(s)
}

func (s *Server) addNewClient(conn net.Conn, name string) *Client {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	c := &Client{
		name: name,
		conn: conn,
	}
	s.peerMap[conn.RemoteAddr()] = c
	return c
}

func (c *Client) readLoop(s *Server) {
	buff := make([]byte, 2048)
	for {
		n, err := c.conn.Read(buff)
		if err != nil {
			log.Println("read err:", err)
			c.conn.Close()
			break
		}
		s.msgch <- Message{
			client:  c,
			payload: buff[:n],
		}

	}
}

type Client struct {
	name string
	conn net.Conn
}

type Message struct {
	client  *Client
	payload []byte
}

func (s *Server) handleMsg() {
	for msg := range s.msgch {
		log.Println(msg.client.name, "wrote: ", string(msg.payload))
		s.broadcastAll(string(msg.client.name) + ": " + string(msg.payload))
	}
}
func (s *Server) broadcastAll(msg string) {
	for _, client := range s.peerMap {
		client.conn.Write([]byte(msg))
	}
}

func main() {
	server := newServer(":3000")
	go server.handleMsg()
	log.Println("Listening...")
	log.Fatal((server.Run()))
}
