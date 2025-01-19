package main

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
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
		go fs.sendFile(3000, conn)
	}
}

func (fs *FileServer) sendFile(size int, conn net.Conn) error {
	file := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, file)
	if err != nil {
		return err
	}

	binary.Write(conn, binary.LittleEndian, int64(size))
	n, err := io.CopyN(conn, bytes.NewReader(file), int64(size))
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
