package main

import (
	"bufio"
	"fmt"
	"io"
	. "mandelbrot/mandelbrot"
	"net"
	"os"
	"strconv"
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
			// Collect parameters
			var xmin, xmax, ymin, ymax float32

			// Helper function to read and parse a float32 value
			readFloat := func(prompt string) (float32, error) {
				writer.WriteString(prompt)
				writer.Flush()
				input, err := reader.ReadString('\n')
				if err != nil {
					return 0, err
				}
				input = strings.TrimSpace(input)
				val, err := strconv.ParseFloat(input, 32)
				if err != nil {
					return 0, fmt.Errorf("invalid float value: %s", input)
				}
				return float32(val), nil
			}

			// Read and parse each parameter
			xmin, err = readFloat("Enter Xmin: \n")
			if err != nil {
				writer.WriteString("Invalid input for Xmin. Please try again.\n")
				writer.Flush()
				continue
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
				fmt.Println("Error generating Mandelbrot image:", err)
				return
			}

			// Save the image with a name based on dimensions
			fileName := fmt.Sprintf("Mandelbrot.png")
			err = mandelbrot.SaveImage(fileName)
			if err != nil {
				fmt.Println("Error saving image:", err)
			} else {
				fmt.Printf("Mandelbrot image saved!\n")
			}

			writer.WriteString("Image generation triggered successfully.\n")
			writer.Flush()

		} else {
			writer.WriteString("Unknown command. Try again.\n")
			writer.Flush()
		}
		err = sendImage(writer)
		if err != nil {
			fmt.Println("Error sending image:", err)
			return
		} else {
			fmt.Println("Image sent successfully.")
		}

	}
}

func sendImage(writer *bufio.Writer) error {
	imageFile, err := os.Open("Mandelbrot.png") // Replace with the path to your PNG file
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
	 Delete the image file after sending it
	err = os.Remove("Mandelbrot.png")
	if err != nil {
		return fmt.Errorf("failed to delete image file: %w", err)
	}

	return nil
}
