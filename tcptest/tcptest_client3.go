package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	serverAddr := "172.21.29.196:8080"

	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Fatal("Failed to connect to server:", err)
	}
	defer conn.Close()
	fmt.Println("Connected to server", serverAddr)

	// Start a goroutine to continuously read from the server
	go readFromServer(conn)

	// Continuously read user input from the terminal and send it to the server
	writeToServer(conn)
}

// Function to read messages from the server continuously
func readFromServer(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		// Read the incoming data line by line
		message, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				// Server closed the connection
				fmt.Println("Server closed the connection.")
			} else {
				// Handle other errors
				fmt.Println("Error reading from server:", err)
			}
			break
		}

		// Print the received message from the server
		fmt.Println("Server says:", message)
	}
}

// Function to send messages to the server based on user input
func writeToServer(conn net.Conn) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		// Read user input from the terminal
		fmt.Print("Enter message to send to server: ")
		scanner.Scan()
		userInput := scanner.Text()

		// Send the input to the server
		if strings.TrimSpace(userInput) != "" {
			_, err := conn.Write([]byte(userInput + "\n"))
			if err != nil {
				fmt.Println("Error sending data to server:", err)
				break
			}
		}
	}
}
