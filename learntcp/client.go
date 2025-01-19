package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func readFromServer() error {
	conn, err := net.Dial("tcp", ":3000")
	if err != nil {
		return err
	}

	file, _ := os.Create("mandelbrot_.png")
	defer conn.Close()
	buf := new(bytes.Buffer)
	for {
		var size int64
		err := binary.Read(conn, binary.LittleEndian, &size)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed by server.")
				return nil
			}
			log.Fatal(err)
		}
		n, err := io.CopyN(file, conn, size)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(buf.Bytes())
		fmt.Printf("recieved %d bytes over the network \n", n)
		buf.Reset()
	}

}

func main() {
	readFromServer()
}
