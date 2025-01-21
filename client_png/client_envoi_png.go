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
	serverAddr := "localhost:8080"

	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Fatal("Error connecting to server:", err)
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
				log.Fatal("Server closed the connection.")
			}
			log.Fatal("Error reading from server:", err)
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
		log.Fatal("Error creating file:", err)
	}
	defer file.Close()

	// Read binary data from the server and save it to the file
	_, err = io.Copy(file, reader)
	if err != nil {
		log.Fatal("Error saving image:", err)
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
				log.Fatal("Error sending to server:", err)
			}
		}

		if userInput == "end" {
			return
		}
	}
}
