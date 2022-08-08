package main

type Queue[T any] []T

func (self *Queue[T]) Push(newQueueItem interface{}) {
	*self = append(*self, newQueueItem)
}

func (self *Queue[T]) Length() int {
	return len(*self)
}

func (self *Queue[T]) Pop() T {
	selfRef := *self
	var element interface{}
	length := selfRef.Length()
	element, *self = selfRef[0], selfRef[1:length]
	return element
}

func NewQueue[T any]() *Queue[T] {
	return &Queue[T]{}
}
