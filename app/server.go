package main

import (
	"fmt"
)

func main() {

	fmt.Println("Starting Go-Redis")

	instance := NewRedisServer()

	go instance.Listen("0.0.0.0:6379")
	go instance.ProcessEventLoop()

	fmt.Println("Server listening, press enter to stop")
	fmt.Scanln()
}