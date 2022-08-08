package main

type Queue []string

func (self *Queue) Push(newQueueItem string) {
	*self = append(*self, newQueueItem)
}

func (self *Queue) Length() int {
	return len(*self)
}

func (self *Queue) Pop() string {
	selfRef := *self
	var element string
	length := selfRef.Length()
	element, *self = selfRef[0], selfRef[1:length]
	return element
}

func NewQueue[T any]() *Queue {
	return &Queue{}
}
