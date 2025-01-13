package tcp

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":8080") // Listen on port 8080
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer ln.Close()
	fmt.Println("Server listening on port 8080")

	for {
		conn, err := ln.Accept() // Accept a connection
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn) // Handle the connection in a new goroutine
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Create a Mandelbrot object and compute the matrix
	mandelbrot := &Mandelbrot{}
	mandelbrot.PrintMandelbrot()

	// Encode the Mandelbrot matrix into bytes
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(mandelbrot)
	if err != nil {
		fmt.Println("Error encoding data:", err)
		return
	}

	// Send the encoded data to the client
	_, err = conn.Write(buffer.Bytes())
	if err != nil {
		fmt.Println("Error sending data:", err)
		return
	}

	fmt.Println("Mandelbrot data sent to client")
}
