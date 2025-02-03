package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"log"
	. "mandelbrot/mandelbrot"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8080")
	//creates TCP server that listens for incoming connections on port8080

	if err != nil {
		log.Fatal("Error starting server:", err)
		//logs the error and stops the program immediately
	}

	defer listener.Close()
	//ensures when main exits the server properly closes
	fmt.Println("Server is listening on port 8080...")

	var wg sync.WaitGroup
	//variable de type sync.WaitGroup => permet aux goroutines de terminer leur execution avant la fin du programme

	for {
		//loops continuously
		conn, err := listener.Accept()
		//listener.Accept() blocks until a client accepts then returns an net.conn object
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			break

		}

		wg.Add(1)                                                            //increment the counter to track an additional goroutine
		fmt.Printf("New client connected: %s\n", conn.RemoteAddr().String()) //prints the client's IP address
		go handleConnection(conn, &wg)                                       // launches a goroutine to handle the client without blocking the server
	}

	wg.Wait()
	//waits for all goroutines to finish before exiting
	fmt.Println("Server shutting down.")
}

func handleConnection(conn net.Conn, wg *sync.WaitGroup) {
	//net.Conn objects represent connection to the client, wg is a pointer to a waitgroup so it can change it's values
	defer wg.Done()
	//Ensures that when the function exits, wg.Done() is called to signal that this goroutine is finished.
	defer conn.Close()
	//Ensures that the client connection is properly closed when the function exits.

	reader := bufio.NewReader(conn)
	//Wraps conn in a buffered reader, making it efficient for reading commands from the client (instead of reading byte by byte)
	writer := bufio.NewWriter(conn)
	//Wraps conn in a buffered writer, allowing efficient writing before flushing data to the client.
	for {
		writer.WriteString("Enter a command (type 'end' to quit, 'send image' to get the image): \n")
		writer.Flush()
		//sends prompt to the client

		command, err := reader.ReadString('\n') //reads input until a newline is received

		if err != nil {
			fmt.Print("Error reading from client:", err)
			return
		}

		command = strings.TrimSpace(command) //removes leading and trailing spaces

		if command == "end" {
			fmt.Print("Client disconnected.")
			return
			// if the user sent "end", the server disconnects the client
		} else if command == "send image" {
			// Collect parameters
			var xmin, xmax, ymin, ymax float32

			// Helper inline function to read and parse a float32 value
			readFloat := func(prompt string) (float32, error) {
				writer.WriteString(prompt)
				writer.Flush()
				input, err := reader.ReadString('\n')
				if err != nil {
					return 0, err
				}
				input = strings.TrimSpace(input)
				val, err := strconv.ParseFloat(input, 32) //converts to float32 precision but still returns float64
				if err != nil {
					return 0, fmt.Errorf("invalid float value: %s", input)
				}
				return float32(val), nil //converts to float64
			}

			// Read and parse each parameter
			xmin, err = readFloat("Enter Xmin: \n")
			if err != nil {
				writer.WriteString("Invalid input for Xmin. Please try again.\n")
				writer.Flush()
				continue //continue makes it ask for the var again
			}

			xmax, err = readFloat("Enter Xmax: \n")
			if err != nil {
				writer.WriteString("Invalid input for Xmax. Please try again.\n")
				writer.Flush()
				continue
			}

			ymin, err = readFloat("Enter Ymin: \n")
			if err != nil {
				writer.WriteString("Invalid input for Ymin. Please try again.\n")
				writer.Flush()
				continue
			}

			ymax, err = readFloat("Enter Ymax: \n")
			if err != nil {
				writer.WriteString("Invalid input for Ymax. Please try again.\n")
				writer.Flush()
				continue
			}

			// Call the mandelbrot function
			writer.WriteString(fmt.Sprintf("generating mandelbrot with xmin=%.2f, xmax=%.2f, ymin=%.2f, ymax=%.2f\n", xmin, xmax, ymin, ymax))
			writer.Flush()

			// Define image dimensions
			const width, height = 1000, 1000
			const numGoRoutines = 100
			const nbIteration = 1000

			mandelbrot := NewMandelbrot(width, height)
			mandelbrot.XMin = float64(xmin)
			mandelbrot.XMax = float64(xmax)
			mandelbrot.YMin = float64(ymin)
			mandelbrot.YMax = float64(ymax)

			err := mandelbrot.PrintOnImage(numGoRoutines, nbIteration)

			if err != nil {
				fmt.Print("Error generating Mandelbrot image:", err)
				return
			}

			fileName := "Mandelbrot.png"
			err = mandelbrot.SaveImage(fileName)
			if err != nil {
				fmt.Println("Error saving image:", err)
			} else {
				fmt.Printf("Mandelbrot image saved!\n")
			}

			writer.WriteString("Image generation triggered successfully.\n")
			writer.Flush()
			err = sendImage(writer)
			if err != nil {
				fmt.Print("Error sending image:", err)
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
	// Read the image file
	imageData, err := os.ReadFile("Mandelbrot.png") //reads the imagefile into imageData
	if err != nil {
		return fmt.Errorf("failed to read image file: %w", err)
	}

	base64Data := base64.StdEncoding.EncodeToString(imageData) //converts image into a base64 string which is easier to transmit using tcp

	// Send the image size first
	sizeMsg := fmt.Sprintf("IMAGE_SIZE:%d\n", len(base64Data))
	_, err = writer.WriteString(sizeMsg)
	if err != nil {
		return fmt.Errorf("failed to send size: %w", err)
	}
	writer.Flush()

	// Send the actual image data
	_, err = writer.WriteString("START_IMAGE\n")
	if err != nil {
		return fmt.Errorf("failed to send start marker: %w", err)
	}
	writer.Flush()

	// Send the base64 data in chunks
	chunkSize := 1024
	for i := 0; i < len(base64Data); i += chunkSize {
		end := i + chunkSize
		if end > len(base64Data) {
			end = len(base64Data)
		}

		_, err := writer.WriteString(base64Data[i:end])
		if err != nil {
			return fmt.Errorf("failed to send data chunk: %w", err)
		}
		writer.Flush()
	}

	// Send end marker
	_, err = writer.WriteString("\nEND_IMAGE\n")
	if err != nil {
		return fmt.Errorf("failed to send end marker: %w", err)
	}
	writer.Flush()

	// Delete the image file
	err = os.Remove("Mandelbrot.png")
	if err != nil {
		return fmt.Errorf("failed to delete image file: %w", err)
	}

	return nil
}
