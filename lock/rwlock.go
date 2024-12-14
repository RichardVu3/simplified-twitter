package lock

import (
	"sync"
)

const MAX_READERS = 32

type RWLock interface {
	Lock()
	RLock()
	RUnlock()
	Unlock()
}

type RWMutex struct {
	mu sync.Mutex
	cond *sync.Cond
	readers int
	writer bool
	waitingWriters int
}

func (rw *RWMutex) Lock() {
	rw.mu.Lock()
	rw.waitingWriters++
	for rw.writer || rw.readers > 0 {
		rw.cond.Wait()
	}
	rw.waitingWriters--
	rw.writer = true
	rw.mu.Unlock()
}

func (rw *RWMutex) RLock() {
	rw.mu.Lock()
	for rw.writer || rw.readers == MAX_READERS || rw.waitingWriters > 0 {
		rw.cond.Wait()
	}
	rw.readers++
	rw.mu.Unlock()
}

func (rw *RWMutex) RUnlock() {
	rw.mu.Lock()
	rw.readers--
	if rw.readers == 0 {
		rw.cond.Broadcast()
	}
	rw.mu.Unlock()
}

func (rw *RWMutex) Unlock() {
	rw.mu.Lock()
	rw.writer = false
	rw.cond.Broadcast()
	rw.mu.Unlock()
}

func NewRWMutex() RWLock {
	rw := &RWMutex{mu: sync.Mutex{}, readers: 0, writer: false, waitingWriters: 0}
	rw.cond = sync.NewCond(&rw.mu)
	return rw
}