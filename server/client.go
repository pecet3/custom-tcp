package main

import (
	"log"
	"net"
)

type Client struct {
	name string
	conn net.Conn
}

func (s *Server) handleClient(conn net.Conn) {
	buff := make([]byte, 255)
	n, err := conn.Read(buff)
	if err != nil {
		log.Println("read nickname err:", err)
		conn.Close()
	}
	name := string(buff[:n])
	log.Println(name)
	s.cchan <- &Client{
		name: name,
		conn: conn,
	}

	log.Println(s.peerMap)
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

	}
}

func (s *Server) addClient(conn net.Conn, name string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.peerMap[name] = &Client{
		name: name,
		conn: conn,
	}

	log.Println(s.peerMap)

}
