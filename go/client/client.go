package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {

	serverAddr := "localhost:8080"
	// Dial connection = initiates a connection can send data !
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		fmt.Print("Error connecting to server:", err)
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
				fmt.Print("Server closed the connection.")
				return
			}
			fmt.Print("Error reading from server:", err)
			return
		}
		message = strings.TrimSpace(message)

		if strings.HasPrefix(message, "IMAGE_SIZE:") { //checks if message string starts with "IMAGE_SIZE:"
			size, err := strconv.Atoi(strings.TrimPrefix(message, "IMAGE_SIZE:")) //trims the prefix from the string and converts it to an integer
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
	//strings.Builder is a buffer used to efficiently build strings.
	//It allows you to append strings to it without creating a new string every time.

	base64Data.Grow(size) // pre-allocates memory in the builder for the expected size

	for base64Data.Len() < size {
		chunk, err := reader.ReadString('\n') // because end message sent is \nEND_IMAGE\n
		if err != nil && err != io.EOF {
			fmt.Println("Error reading image data:", err)
			return
		}

		// Check if we've reached the end marker
		//if strings.TrimSpace(chunk) == "END_IMAGE" {
		//	break
		//}

		base64Data.WriteString(strings.TrimSpace(chunk)) //chunk is added to the base64Data string builder
		fmt.Printf("Received %d/%d bytes\n", base64Data.Len(), size)
	}

	// Decode base64 data
	imageData, err := base64.StdEncoding.DecodeString(base64Data.String())
	if err != nil {
		//log.Fatal("Error creating file:", err)
		fmt.Println("Error decoding image data and creating file:", err)
		return
	}

	// Save the image
	err = os.WriteFile("received_image.png", imageData, 0644) // writes the decoded image data to a file, 0644 specifies that the file will be readable and writable by the owner, and readable by others.
	if err != nil {
		fmt.Print("Error saving image:", err)
		return
	}

	fmt.Println("Image received and saved as 'received_image.png'")
}

func writeToServer(conn net.Conn) {
	scanner := bufio.NewScanner(os.Stdin) //Creates a new Scanner object that reads input from the standard input
	for {
		fmt.Print("Enter command: ")
		scanner.Scan()              //reads from the console until enter
		userInput := scanner.Text() // retrieves as a string
		if strings.TrimSpace(userInput) != "" {
			_, err := conn.Write([]byte(userInput + "\n")) //converts the string to a byte slice required by conn.Write
			if err != nil {
				fmt.Print("Error sending to server:", err)
				return
			}
		}
		if userInput == "end" {
			return
		}
	}
}
