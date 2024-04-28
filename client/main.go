package main

import (
	"fmt"
	"net"
)

func main() {
	serverAddr := "localhost:3000"

	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Connection error:", err)
		return
	}
	defer conn.Close()

	var name string
	for {
		fmt.Print("Enter your name: ")
		_, err = fmt.Scan(&name)
		if err != nil {
			fmt.Println(err)
		}
		if len(name) < 2 {
			fmt.Println("You must enter at least name with 3 characters")
			continue
		}
		handleHandShake(conn, name)
		break
	}

	go readLoop(conn)
	for {
		var msg string
		_, err = fmt.Scan(&msg)

		_, err := conn.Write([]byte(msg))
		if err != nil {
			fmt.Println("Error message write:", err)
			return
		}
	}
}

func handleHandShake(conn net.Conn, name string) {
	buff := make([]byte, 255)
	n, err := conn.Read(buff)
	if err != nil {
		fmt.Println("Handshake response error:", err)
		return
	}
	response := string(buff[:n])
	if response == "ENTER_NAME"+conn.LocalAddr().String() {
		conn.Write([]byte(name))
	}
}

func readLoop(conn net.Conn) {
	response := make([]byte, 2048)

	for {
		n, err := conn.Read(response)
		if err != nil {
			fmt.Println("Error message response:", err)
			return
		}

		fmt.Println(string(response[:n]))
	}

}
