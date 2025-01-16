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
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Server is listening on port 8080...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		fmt.Printf("New client connected: %s\n", conn.RemoteAddr().String())
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		writer.WriteString("Enter a command (type 'end' to quit, 'send image' to get the image): \n")
		writer.Flush()

		command, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from client:", err)
			return
		}

		command = strings.TrimSpace(command)

		if command == "end" {
			fmt.Println("Client disconnected.")
			return
		} else if command == "send image" {
			err := sendImage(writer)
			if err != nil {
				fmt.Println("Error sending image:", err)
				return
			}
			fmt.Println("Image sent successfully.")
		} else {
			writer.WriteString("Unknown command. Try again.\n")
			writer.Flush()
		}
	}
}

func sendImage(writer *bufio.Writer) error {
	imageFile, err := os.Open("nolachacha.png") // Replace with the path to your PNG file
	if err != nil {
		return fmt.Errorf("failed to open image file: %w", err)
	}
	defer imageFile.Close()

	// Inform the client that binary data is being sent
	writer.WriteString("START_IMAGE\n")
	writer.Flush()

	// Send the image file as binary data
	_, err = io.Copy(writer, imageFile)
	if err != nil {
		return fmt.Errorf("failed to send image file: %w", err)
	}

	writer.Flush()
	return nil
}
