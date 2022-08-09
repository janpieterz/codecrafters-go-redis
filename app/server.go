package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
	// Uncomment this block to pass the first stage
	// "net"
	// "os"
)

var memoryCache MemoryCache

func main() {
	memoryCache = make(MemoryCache)
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
			ParseInput(nextItem)
		} else {
			time.Sleep(1 * time.Millisecond)
		}
	}
}

func ParseInput(message RedisMessage) {
	fmt.Printf("Processing messages %s \n", message.messages)
	if len(message.messages) < 1 {
		fmt.Println("Something went wrong, no messages in input message")
	}

	switch message.messages[0] {
	case "echo":
		fmt.Println("Sending back echo response")
		SendMessage(message.messages[1], message.connection)
		break
	case "ping":
		fmt.Println("Sending back PONG response")
		SendMessage("PONG", message.connection)
		break
	case "set":
		fmt.Println("Setting value")
		SetValue(message.messages[1:])

		SendMessage("OK", message.connection)
		break
	case "get":
		fmt.Println("Getting value")
		fmt.Println("Cache status:", memoryCache)
		value := GetValue(message.messages[1])
		SendMessage(value, message.connection)
		break
	}
}

func SetValue(parameters []string) {
	var expirationTime *time.Time
	expirationTime = nil
	if len(parameters) > 2 {
		if parameters[2] == "px" {
			input, error := strconv.Atoi(parameters[3])
			if error != nil {
				fmt.Printf("Could not parse input, error: %s \n", error.Error())
			} else {
				expiration := time.Now().Add(time.Duration(input * 1_000_000))
				expirationTime = &expiration
			}
		}
	}
	fmt.Printf("Setting key %s with value %s and expiration %d \n", parameters[0], parameters[1], expirationTime)
	memoryCache.Push(parameters[0], parameters[1], expirationTime)
}

func GetValue(key string) string {
	value := memoryCache.Get(key)
	return value
}

func SendMessage(message string, connection net.Conn) {
	formattedMessage := message
	if message != "$-1" {
		formattedMessage = formatRESPString(message)
	}
	_, err := connection.Write([]byte(formattedMessage))
	if err != nil {
		fmt.Println("Error sending data", err.Error())
		os.Exit(1)
	}
}

func formatRESPString(input string) string {
	return "+" + input + "\r\n"
}
