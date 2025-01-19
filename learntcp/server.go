package main

import (
	"fmt"
	"io"
	"log"
	. "mandelbrot/mandelbrot"
	"net"
	"os"
)

type FileServer struct{}

func (fs *FileServer) start() {

	// ln is a TCP net.Listener it can listen to the connection of port 3000
	ln, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go fs.sendFile(1000, 1000, conn)
	}
}

func (fs *FileServer) sendFile(width uint32, height uint32, conn net.Conn) error {
	mandelbrot := NewMandelbrot(width, height)

	err := mandelbrot.PrintOnImage(10, 100)

	if err != nil {
		fmt.Println("Error generating Mandelbrot image:", err)
		return err
	}

	mandelbrot.SaveImage("mandelbrot.png")
	size := width * height
	/*
		file := make([]byte, size)
		a, err := io.ReadFull(rand.Reader, file)
		if err != nil {
			return err
		}

		binary.Write(conn, binary.LittleEndian, int64(size))
		n, err := io.CopyN(conn, bytes.NewReader(file), int64(size))
	*/
	file, err := os.Open("mandelbrot.png")
	n, err := io.CopyN(conn, file, int64(size))
	if err != nil {
		return err
	}
	fmt.Printf("written %d bytes over the network", n)
	return nil
}

func main() {
	server := &FileServer{}
	server.start()
}
