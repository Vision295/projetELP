package main

import (
	"bufio"
	"fmt"
	"image/png"
	"io"
	. "mandelbrot/mandelbrot"
	"net"
	"os"
	"strconv"
)

func main() {
	// Create a TCP server that listens on port 8080
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
	writer.Write([]byte("Enter the size of the image Mandelbrot you want to create. \nFirst enter its width : \n"))
	writer.Flush() // Flush to send the message immediately
	strwidth, _ := bufio.NewReader(conn).ReadString('\n')
	// Convert the input strings to integers
	width, err1 := strconv.Atoi(strwidth)

	writer.Write([]byte("\nThen, enter its height : \n"))
	writer.Flush() // Flush to send the message immediately
	strheight, _ := bufio.NewReader(conn).ReadString('\n')
	print(strheight)
	// Convert the input strings to integers
	height, err2 := strconv.Atoi(strheight)

	if err1 != nil || err2 != nil {
		writer.Write([]byte("Error: Both inputs must be integers.\n"))
		writer.Flush() // Flush to send the error message immediately
		return
	}

	var mandelbrot = NewMandelbrot()
	mandelbrot.Width = uint32(width)
	mandelbrot.Height = uint32(height)

	mandelbrot.PrintOnImage(100)
	file, _ := os.Create("mandelbrot.png")
	defer file.Close()
	png.Encode(file, mandelbrot.Image)

	io.Copy(conn, file)

	// Send the result back to the client
	writer.Write([]byte(fmt.Sprintf("The sum is\n")))
	writer.Flush() // Flush to send the result immediately
}
