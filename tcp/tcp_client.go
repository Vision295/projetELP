package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	. "mandelbrot/mandelbrot"
	"net"
)

func ConnectToServer(ip string) {
	conn, err := net.Dial("tcp6", "["+ip+"]"+":8080") // Replace with your server's IP and port
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	// Read the response from the server
	var buffer bytes.Buffer
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			break
		}
		buffer.Write(buf[:n])
	}

	// Decode the received data
	var mandelbrot Mandelbrot
	decoder := gob.NewDecoder(&buffer)
	err = decoder.Decode(&mandelbrot)
	if err != nil {
		fmt.Println("Error decoding data:", err)
		return
	}

	// Use the Mandelbrot data (e.g., display or process it)
	fmt.Println("Received Mandelbrot matrix:", mandelbrot)
}

func main() {

	ConnectToServer("fe80::215:5dff:fe77:9ba0/64")
}
