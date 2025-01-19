package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
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
	// Read the "START_IMAGE" line with the size
	header, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading image header:", err)
		return
	}

	// Parse the file size
	parts := strings.Fields(header)
	if len(parts) != 2 || parts[0] != "START_IMAGE" {
		fmt.Println("Invalid START_IMAGE header:", header)
		return
	}

	fileSize, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		fmt.Println("Invalid file size:", parts[1])
		return
	}

	// Create the output file
	file, err := os.Create("received_image.png")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Read exactly `fileSize` bytes for the image
	_, err = io.CopyN(file, reader, fileSize)
	if err != nil {
		fmt.Println("Error saving image:", err)
		return
	}

	// Read the "END_IMAGE" marker (optional if not used)
	endMarker, _ := reader.ReadString('\n')
	if strings.TrimSpace(endMarker) != "END_IMAGE" {
		fmt.Println("Warning: Missing END_IMAGE marker.")
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
