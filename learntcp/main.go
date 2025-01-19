package main

import (
	"bytes"
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
		go fs.readLoop(conn)
	}
}

func (fs *FileServer) readLoop(conn net.Conn) {
	defer conn.Close()
	buf := new(bytes.Buffer)
	for {
		var size int64
		err := binary.Read(conn, binary.LittleEndian, &size)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed by client.")
				return
			}
			log.Fatal(err)
		}
		n, err := io.CopyN(buf, conn, size)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(buf.Bytes())
		fmt.Printf("recieved %d bytes over the network \n", n)
		buf.Reset()
	}
}

func main() {

	server := &FileServer{}
	server.start()
}
