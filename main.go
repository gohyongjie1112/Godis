package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

const Port = ":6379"

func main() {

	listener, err := net.Listen("tcp", Port)
	if err != nil {
		fmt.Printf("Error listening on port %s: %v\n", Port, err)
	}
	connection, err := listener.Accept()
	if err != nil {
		fmt.Printf("Error accepting connection: %v\n", err)
	}

	defer connection.Close()

	for {
		buffer := make([]byte, 1024)
		length, err := connection.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Error reading from client: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Received %d bytes: %s\n", length, string(buffer[:length]))
		connection.Write([]byte("+PONG\r\n"))
	}

}
