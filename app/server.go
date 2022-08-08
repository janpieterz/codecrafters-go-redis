package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
	// Uncomment this block to pass the first stage
	// "net"
	// "os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Starting Go-Redis")

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

	eventLoop := NewQueue()
	buffer := make([]byte, 256)

	go ProcessEventLoop(eventLoop, &connection)

	for {
		receivedCount, err := connection.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Found EOF")
			} else {
				fmt.Println("Error reading input", err.Error())
				os.Exit(1)
			}
		}
		if receivedCount == 0 {
			continue
		}
		input := string(buffer[:receivedCount])
		splits := strings.Split(input, "\r\n")
		for _, split := range splits {
			fmt.Println("Received '%s', adding to event loop", split)
			eventLoop.Push(split)
		}
	}
}

func ProcessEventLoop(queue *Queue, connection *net.Conn) {
	for {
		if queue.Length() > 0 {
			nextItem := queue.Pop()
			parseInput(nextItem, *connection)
		} else {
			time.Sleep(1 * time.Millisecond)
		}
	}
}

func parseInput(input string, connection net.Conn) {
	fmt.Printf("Processing %s \n", input)
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
