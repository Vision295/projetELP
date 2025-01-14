package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	// Define the server address (can be a TCP address or Unix socket)
	serverAddr := "localhost:8080"

	// Dial the server (connect to it)
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Fatal("Failed to connect to server:", err)
	}
	defer conn.Close()

	// Print a message indicating that we're connected
	fmt.Println("Connected to server", serverAddr)

	// Start a goroutine to continuously read from the server
	go readFromServer(conn)

	// Continuously read user input from the terminal and send it to the server
	writeToServer(conn)
	// erreur car il faut attendre le résultat de go readFromServer
	outFile, err := os.Create("OutputFile.png")
	if err != nil {
		log.Fatal("Error creating files", err)
	}

	defer outFile.Close()

	i, err := io.Copy(outFile, conn)
	if err != nil {
		log.Fatal("Error initialing files ", err)
	}
	fmt.Print(i)
	fmt.Print("ok")
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
