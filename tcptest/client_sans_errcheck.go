package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	serverAddr := "172.21.29.196:8080"

	conn, _ := net.Dial("tcp", serverAddr)
	defer conn.Close()
	fmt.Println("Connected to server", serverAddr)

	go readFromServer(conn)

	writeToServer(conn)
}

func readFromServer(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				// Server closed the connection
				fmt.Println("Server closed the connection.")
				return
			} else {
				// Handle other errors
				fmt.Println("Error reading from server:", err)
			}
			break
		}
		fmt.Println(message, err)
	}
}

func writeToServer(conn net.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter message to send to server: ")
		scanner.Scan()
		userInput := scanner.Text()
		if strings.TrimSpace(userInput) != "" {
			conn.Write([]byte(userInput + "\n"))
		}
	}
}
