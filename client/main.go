package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// Adres serwera TCP
	serverAddr := "localhost:3000"

	// Nawiązanie połączenia z serwerem
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Błąd połączenia:", err)
		return
	}
	defer conn.Close()

	// Pętla nieskończona do wysyłania i odbierania wiadomości
	for {
		// Odczytanie wiadomości od użytkownika
		fmt.Print("Wpisz wiadomość: ")
		reader := bufio.NewReader(os.Stdin)
		msg, _ := reader.ReadString('\n')

		// Wysłanie wiadomości do serwera
		_, err := conn.Write([]byte(msg))
		if err != nil {
			fmt.Println("Błąd wysyłania wiadomości:", err)
			return
		}

		// Odbieranie odpowiedzi od serwera
		response := make([]byte, 2048)
		n, err := conn.Read(response)
		if err != nil {
			fmt.Println("Błąd odbierania wiadomości:", err)
			return
		}

		// Wyświetlenie otrzymanej odpowiedzi od serwera
		fmt.Println("Odpowiedź serwera:", string(response[:n]))
	}
}
