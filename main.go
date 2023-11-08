package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

const Port = ":6379"

// RESP string format is: $<number of bytes>\r\n<bytes>\r\n
// Example: $6\r\nfoobar\r\n
// https://redis.io/topics/protocol
func readInput(input string) {
	reader := bufio.NewReader(strings.NewReader(input))
	b, err := reader.ReadByte()
	if err != nil {
		fmt.Printf("Error reading byte: %v\n", err)
		os.Exit(1)
	}
	if b != '$' {
		fmt.Println("Invalid string format, expected $ but got", string(b))
		os.Exit(1)
	}
	size, err := reader.ReadByte()
	if err != nil {
		fmt.Printf("Error reading byte: %v\n", err)
		os.Exit(1)
	}
	strSize, err := strconv.ParseInt(string(size), 10, 64)
	if err != nil {
		fmt.Printf("Error converting to int: %v\n", err)
		os.Exit(1)
	}

	// consume /r/n
	reader.ReadByte()
	reader.ReadByte()

	name := make([]byte, strSize)
	reader.Read(name)

	fmt.Println(string(name))
}

func main() {

	readInput("$6\r\nfoobar\r\n")

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
