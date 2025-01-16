package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Server is listening on port 8080...")

	// Handle incoming connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		fmt.Printf("New client connected: %s\n", conn.RemoteAddr().String())

		// Handle the connection in a separate goroutine
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Create a buffered writer
	writer := bufio.NewWriter(conn)

	// Ask the client for the first number
	writer.Write([]byte("Please enter the first number: \n"))
	writer.Flush() // Flush to send the message immediately

	num1Str, _ := bufio.NewReader(conn).ReadString('\n')
	num1Str = strings.TrimSpace(num1Str)

	// Ask the client for the second number
	writer.Write([]byte("Please enter the second number: \n"))
	writer.Flush() // Flush to send the message immediately

	num2Str, _ := bufio.NewReader(conn).ReadString('\n')
	num2Str = strings.TrimSpace(num2Str)

	// Convert the input strings to integers
	num1, err1 := strconv.Atoi(num1Str)
	num2, err2 := strconv.Atoi(num2Str)

	if err1 != nil || err2 != nil {
		writer.Write([]byte("Error: Both inputs must be integers.\n"))
		writer.Flush() // Flush to send the error message immediately
		return
	}

	// Calculate the sum
	sum := num1 + num2

	// Send the result back to the client
	writer.Write([]byte(fmt.Sprintf("The sum is: %d\n", sum)))
	writer.Flush() // Flush to send the result immediately
}
