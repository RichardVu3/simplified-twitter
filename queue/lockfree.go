package queue

import (
	"sync/atomic"
)

type Request struct {
	ID        int     `json:"id"`
	Command   string  `json:"command"`
	Body      string  `json:"body"`
	Timestamp float64 `json:"timestamp"`
}

func NewRequest() *Request {
	return &Request{}
}

type Node struct {
	request *Request
	next atomic.Pointer[Node]
}

func NewNode(request *Request) *Node {
	return &Node{request: request}
}

type LockFreeQueue struct {
	head atomic.Pointer[Node]
	tail atomic.Pointer[Node]
}

func NewLockFreeQueue() *LockFreeQueue {
	var head atomic.Pointer[Node]
	var tail atomic.Pointer[Node]
	node := Node{}
	head.Store(&node)
	tail.Store(&node)
	newQueue := &LockFreeQueue{}
	newQueue.head.Store(head.Load())
	newQueue.tail.Store(tail.Load())
	return newQueue
}

func (queue *LockFreeQueue) Enqueue(task *Request) {
	newNode := NewNode(task)
	for {
		tail := queue.tail.Load()
		next := tail.next.Load()
		if tail == queue.tail.Load() {
			if next == nil {
				if tail.next.CompareAndSwap(next, newNode) {
					queue.tail.CompareAndSwap(tail, newNode)
					return
				}
			} else {
				queue.tail.CompareAndSwap(tail, next)
			}
		}
	}
}

func (queue *LockFreeQueue) Dequeue() *Request {
	var req *Request
	for {
		head := queue.head.Load()
		tail := queue.tail.Load()
		first := head.next.Load()
		if head == queue.head.Load() {
			if head == tail {
				if first == nil {
					return req
				}
				queue.tail.CompareAndSwap(tail, first)
			} else {
				request := first.request
				if queue.head.CompareAndSwap(head, first) {
					return request
				}
			}
		}
	}
}

func (queue *LockFreeQueue) IsEmpty() bool {
	head := queue.head.Load()
	tail := queue.tail.Load()
	return head == tail && head.next.Load() == nil
}
