package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
)

var err error

func main() {
	serverAddr := "172.21.29.196:8080"

	conn, _ := net.Dial("tcp", serverAddr)
	defer conn.Close()
	fmt.Println("Connected to server", serverAddr)
	isConnectionUp := false
	if err != nil {
		isConnectionUp = true
	}
	go readFromServer(conn, isConnectionUp)

	err = writeToServer(conn)
}

func readFromServer(conn net.Conn, isConnectionUp bool) {
	reader := bufio.NewReader(conn)
	for isConnectionUp {
		message, err := reader.ReadString('\n')
		fmt.Println(message, err)
	}
}

func writeToServer(conn net.Conn) error {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Enter message to send to server: ")
		scanner.Scan()
		userInput := scanner.Text()
		if strings.TrimSpace(userInput) != "" {
			conn.Write([]byte(userInput + "\n"))
		}
		if userInput == "end" {
			return errors.New("end")
		}
	}
}
