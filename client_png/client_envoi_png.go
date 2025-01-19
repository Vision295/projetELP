package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func main() {
	serverAddr := "localhost:8080"

	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
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
			if err == io.EOF {
				fmt.Println("Server closed the connection.")
				return
			}
			fmt.Println("Error reading from server:", err)
			return
		}

		message = strings.TrimSpace(message)

		if message == "START_IMAGE" {
			saveImage(reader)
		} else {
			fmt.Println(message)
		}
	}
}

func saveImage(reader *bufio.Reader) {
	file, err := os.Create("received_image.png") // File to save the received image
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Read binary data from the server and save it to the file
	_, err = io.Copy(file, reader)
	if err != nil {
		fmt.Println("Error saving image:", err)
		return
	}

	fmt.Println("Image received and saved as 'received_image.png'.")
}

func writeToServer(conn net.Conn) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Enter command: ")
		scanner.Scan()
		userInput := scanner.Text()

		if strings.TrimSpace(userInput) != "" {
			_, err := conn.Write([]byte(userInput + "\n"))
			if err != nil {
				fmt.Println("Error sending to server:", err)
				return
			}
		}

		if userInput == "end" {
			return
		}
	}
}
