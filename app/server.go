package main

import (
	"fmt"
	"io"
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

	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	connection, err := listener.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}
	fmt.Println("Listening to input")
	defer connection.Close()

	for {
		bytes, err := io.ReadAll(connection)
		if err != nil {
			fmt.Println("Error reading input", err.Error())
			os.Exit(1)
		}
		input := string(bytes)
		parseInput(input, connection)
	}
}

func parseInput(input string, connection net.Conn) {
	fmt.Printf("Received %s", input)
	if input == "ping" {
		_, err := connection.Write(formatRESPString("PONG"))
		if err != nil {
			fmt.Println("Error sending data", err.Error())
			os.Exit(1)
		}
	}
}

func formatRESPString(input string) []byte {
	formattedString := "+" + input + "\r\n"
	return []byte(formattedString)
}
