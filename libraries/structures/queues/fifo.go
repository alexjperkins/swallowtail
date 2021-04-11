package queues

import (
	"fmt"
	"sync"
)

type FIFOQueue struct {
	q       []interface{}
	mtx     sync.RWMutex
	maxSize int
}

func NewFIFOQueue(maxSize int) *FIFOQueue {
	return &FIFOQueue{
		q:       []interface{}{},
		mtx:     sync.RWMutex{},
		maxSize: maxSize,
	}
}

func (f *FIFOQueue) Pop() (interface{}, bool) {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	if len(f.q) == 0 {
		return nil, false
	}
	r := f.q[0]
	f.q = f.q[1:]
	return r, true
}

func (f *FIFOQueue) Peek() (interface{}, bool) {
	f.mtx.RLock()
	defer f.mtx.RUnlock()
	if len(f.q) == 0 {
		return nil, false
	}
	return f.q[len(f.q)-1], true
}

func (f *FIFOQueue) Push(item interface{}) error {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	if len(f.q) == f.maxSize {
		return fmt.Errorf("FIFO already at max size: %d", f.maxSize)
	}
	f.q = append(f.q, item)
	return nil
}

func (f *FIFOQueue) Len() int {
	f.mtx.RLock()
	defer f.mtx.RUnlock()
	return len(f.q)
}

func (f *FIFOQueue) AtCapacity() bool {
	return f.Len() == f.maxSize
}

func (f *FIFOQueue) GetAsArray() []interface{} {
	f.mtx.Lock()
	defer f.mtx.Unlock()
	cp := make([]interface{}, len(f.q))
	copy(cp, f.q)
	return cp
}
