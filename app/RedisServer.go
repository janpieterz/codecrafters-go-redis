package main

import (
	"fmt"
	uuid2 "github.com/google/uuid"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type RedisMessage struct {
	messages []string

	connectionId uuid2.UUID
}

type RedisServer struct {
	MemoryCache MemoryCache
	Connections map[uuid2.UUID]net.Conn
	EventLoop   chan RedisMessage
}

func NewRedisServer() *RedisServer {
	server := RedisServer{make(map[string]CacheItem), make(map[uuid2.UUID]net.Conn), make(chan RedisMessage)}
	return &server
}

func (server *RedisServer) Listen(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	for {
		connection, err := listener.Accept()
		defer connection.Close()

		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		uuid, err := uuid2.NewUUID()
		if err != nil {
			fmt.Println("Error generating unique id: ", err.Error())
			os.Exit(1)
		}
		server.Connections[uuid] = connection

		fmt.Printf("Connected to new client %s, starting listening.\n", uuid.String())
		go server.ListenToConnection(uuid)

	}
}

func (server *RedisServer) ProcessEventLoop() {
	for {
		nextMessage := <-server.EventLoop
		server.ParseInput(nextMessage)
	}
}

func (server *RedisServer) ListenToConnection(connectionId uuid2.UUID) {
	connection := server.Connections[connectionId]
	if connection == nil {
		fmt.Printf("Connection %s does not exist anymore", connectionId.String())
		return
	}

	for {
		buffer := make([]byte, 256)
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
		server.ProcessMessage(input, connectionId)
	}
}

func (server *RedisServer) ProcessMessage(message string, connectionId uuid2.UUID) {
	splits := strings.Split(message, "\r\n")
	if len(splits) < 1 {
		return
	}

	messages := make([]string, 0)

	for _, split := range splits {
		if len(split) < 1 {
			continue
		}
		firstCharacter := split[0]
		if firstCharacter == '*' || firstCharacter == '$' {
			continue
		}
		messages = append(messages, split)
		fmt.Printf("Received '%s', adding to event loop \n", split)
	}
	server.EventLoop <- RedisMessage{messages: messages, connectionId: connectionId}
}

func (server *RedisServer) ParseInput(message RedisMessage) {
	connection := server.Connections[message.connectionId]
	if connection == nil {
		fmt.Printf("Connection %s does not exist anymore", message.connectionId.String())
		return
	}

	fmt.Printf("Processing messages %s \n", message.messages)
	if len(message.messages) < 1 {
		fmt.Println("Something went wrong, no messages in input message")
	}

	switch message.messages[0] {
	case "echo":
		fmt.Println("Sending back echo response")
		SendMessage(message.messages[1], connection)
		break
	case "ping":
		fmt.Println("Sending back PONG response")
		SendMessage("PONG", connection)
		break
	case "set":
		fmt.Println("Setting value")
		server.SetValue(message.messages[1:])

		SendMessage("OK", connection)
		break
	case "get":
		fmt.Println("Getting value")
		fmt.Println("Cache status:", server.MemoryCache)
		value := server.GetValue(message.messages[1])
		SendMessage(value, connection)
		break
	}
}

func (server *RedisServer) SetValue(parameters []string) {
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
	server.MemoryCache.Push(parameters[0], parameters[1], expirationTime)
}

func (server *RedisServer) GetValue(key string) string {
	value := server.MemoryCache.Get(key)
	return value
}

func SendMessage(message string, connection net.Conn) {
	formattedMessage := message
	if message != "nil" {
		formattedMessage = formatRESPString(message)
	} else {
		formattedMessage = "$-1\r\n"
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
