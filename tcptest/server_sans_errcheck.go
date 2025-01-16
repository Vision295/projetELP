package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func main() {
	listener, _ := net.Listen("tcp", ":8080")
	defer listener.Close()
	fmt.Println("Server is listening on port 8080...")

	for {
		conn, _ := listener.Accept()

		fmt.Printf("New client connected: %s\n", conn.RemoteAddr().String())

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	num1Str := ""
	for num1Str != "end" {
		writer := bufio.NewWriter(conn)

		writer.Write([]byte("Please enter what you want : \n"))
		writer.Flush()

		num1Str, _ = bufio.NewReader(conn).ReadString('\n')
		num1Str = strings.TrimSpace(num1Str)
	}
}
