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

	eventLoop := NewQueue()

	go ProcessEventLoop(eventLoop)

	for {
		connection, err := listener.Accept()
		defer connection.Close()

		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		fmt.Println("Connected to new client, start listening")

		go ListenToConnection(connection, eventLoop)

	}
}

func ListenToConnection(connection net.Conn, eventLoop *Queue) {
	buffer := make([]byte, 256)

	for {
		receivedCount, err := connection.Read(buffer)
		if err != nil {
			if err == io.EOF {
				continue
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
		if len(splits) < 1 {
			continue
		}

		messages := make([]string, 0)

		for _, split := range splits {
			if len(split) >= 1 {
				firstCharacter := split[0]
				if firstCharacter == '*' || firstCharacter == '$' {
					continue
				}
			}
			if len(split) < 1 {
				continue
			}
			messages = append(messages, split)
			fmt.Printf("Received '%s', adding to event loop \n", split)
		}
		eventLoop.Push(RedisMessage{messages: messages, connection: connection})
	}
}

func ProcessEventLoop(queue *Queue) {
	for {
		if queue.Length() > 0 {
			nextItem := queue.Pop()
			parseInput(nextItem)
		} else {
			time.Sleep(1 * time.Millisecond)
		}
	}
}

func parseInput(message RedisMessage) {
	fmt.Printf("Processing messages %s \n", message.messages)
	if len(message.messages) < 1 {
		fmt.Println("Something went wrong, no messages in input message")
	}

	switch message.messages[0] {
	case "echo":
		SendMessage(message.messages[1], message.connection)
		break
	case "ping":
		SendMessage("PONG", message.connection)
	}
}

func SendMessage(message string, connection net.Conn) {
	_, err := connection.Write(formatRESPString(message))
	if err != nil {
		fmt.Println("Error sending data", err.Error())
		os.Exit(1)
	}
}

func formatRESPString(input string) []byte {
	formattedString := "+" + input + "\r\n"
	return []byte(formattedString)
}
