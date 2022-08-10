package main

import (
	uuid2 "github.com/google/uuid"
)

type RedisMessage struct {
	messages []string

	connectionId uuid2.UUID
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
