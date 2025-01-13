package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	. "mandelbrot/mandelbrot"
	"net"
)

func HandleConnection(conn net.Conn, mandelbrot Mandelbrot) {
	defer conn.Close()
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

func LaunchServerTCP(mandelbrot Mandelbrot) {
	ln, err := net.Listen("tcp6", "[fe80::215:5dff:fe77:9ba0/64]:8080") // Listen on port 8080
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

		go HandleConnection(conn, mandelbrot) // Handle the connection in a new goroutine
	}
}

var mandelbrot = NewMandelbrot()

func main() {

	LaunchServerTCP(mandelbrot)
}
