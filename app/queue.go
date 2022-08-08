package main

import "net"

type QueuedMessage struct {
	message    string
	connection net.Conn
}

type Queue []QueuedMessage

func (self *Queue) Push(newQueueItem QueuedMessage) {
	*self = append(*self, newQueueItem)
}

func (self *Queue) Length() int {
	return len(*self)
}

func (self *Queue) Pop() QueuedMessage {
	selfRef := *self
	var element QueuedMessage
	length := selfRef.Length()
	element, *self = selfRef[0], selfRef[1:length]
	return element
}

func NewQueue() *Queue {
	return &Queue{}
}
