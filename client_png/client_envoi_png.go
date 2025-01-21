package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
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

		if strings.HasPrefix(message, "IMAGE_SIZE:") {
			size, err := strconv.Atoi(strings.TrimPrefix(message, "IMAGE_SIZE:"))
			if err != nil {
				fmt.Println("Error parsing image size:", err)
				continue
			}
			receiveImage(reader, size)
		} else {
			fmt.Println(message)
		}
	}
}

func receiveImage(reader *bufio.Reader, size int) {
	// Wait for start marker
	marker, err := reader.ReadString('\n')
	if err != nil || strings.TrimSpace(marker) != "START_IMAGE" {
		fmt.Println("Error: Expected START_IMAGE marker")
		return
	}

	// Read the base64 data
	var base64Data strings.Builder
	base64Data.Grow(size) // Pre-allocate the required size

	for base64Data.Len() < size {
		chunk, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			fmt.Println("Error reading image data:", err)
			return
		}

		// Check if we've reached the end marker
		if strings.TrimSpace(chunk) == "END_IMAGE" {
			break
		}

		base64Data.WriteString(strings.TrimSpace(chunk))
		fmt.Printf("Received %d/%d bytes\n", base64Data.Len(), size)
	}

	// Decode base64 data
	imageData, err := base64.StdEncoding.DecodeString(base64Data.String())
	if err != nil {
		log.Fatal("Error creating file:", err)
		fmt.Println("Error decoding image data:", err)
		return
	}

	// Save the image
	err = os.WriteFile("received_image.png", imageData, 0644)
	if err != nil {
		log.Fatal("Error saving image:", err)
	}

	fmt.Println("Image received and saved as 'received_image.png'")
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
