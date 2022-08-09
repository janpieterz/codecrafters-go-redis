package main

import "net"

type RedisMessage struct {
	messages []string

	connection net.Conn
}

type Queue []RedisMessage

func (self *Queue) Push(newQueueItem RedisMessage) {
	*self = append(*self, newQueueItem)
}

func (self *Queue) Length() int {
	return len(*self)
}

func (self *Queue) Pop() RedisMessage {
	selfRef := *self
	var element RedisMessage
	length := selfRef.Length()
	element, *self = selfRef[0], selfRef[1:length]
	return element
}

func NewQueue() *Queue {
	return &Queue{}
}
