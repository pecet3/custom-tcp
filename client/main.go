package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	serverAddr := "localhost:3000"

	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Connection error:", err)
		return
	}
	defer conn.Close()
	handleHandShake(conn)
	for {
		fmt.Print("Message: ")
		reader := bufio.NewReader(os.Stdin)
		msg, _ := reader.ReadString('\n')

		_, err := conn.Write([]byte(msg))
		if err != nil {
			fmt.Println("Error message write:", err)
			return
		}

		response := make([]byte, 2048)
		n, err := conn.Read(response)
		if err != nil {
			fmt.Println("Error message response:", err)
			return
		}

		fmt.Println("[SERVER]:", string(response[:n]))
	}
}

func handleHandShake(conn net.Conn) {
	for {
		buff := make([]byte, 255)
		n, err := conn.Read(buff)
		if err != nil {
			fmt.Println("Handshake response error:", err)
			return
		}
		response := string(buff[:n])
		if response == "name" {
			fmt.Print("Your name: ")
			reader := bufio.NewReader(os.Stdin)
			msg, _ := reader.ReadString('\n')
			conn.Write([]byte(msg))
		}
		fmt.Println("[SERVER]:", string(buff[:n]))
	}
}
