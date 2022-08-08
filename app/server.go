package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	// Uncomment this block to pass the first stage
	// "net"
	// "os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	//Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	connection, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	_, err = connection.Write(formatRESPString("PONG"))
	if err != nil {
		fmt.Println("Error writing to connection", err.Error())
		os.Exit(1)
	}
}

func formatRESPString(input string) []byte {
	formattedString := "+" + input + "\r\n"
	var buf bytes.Buffer
	fmt.Fprintf(&buf, formattedString)
	return buf.Bytes()
}
