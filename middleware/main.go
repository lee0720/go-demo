package main

import (
	"fmt"
	"sync"
	"time"

	"gitlab.com/lilh/go-demo/middleware/queue"
)

const (
	signalInterval = 200
	signalChanSize = 10
)

type SafeQueue struct {
	q *queue.Queue
	sync.Mutex
	C chan struct{}
}

func NewSafeQueue() *SafeQueue {
	sq := &SafeQueue{
		q: queue.New(),
		C: make(chan struct{}, signalChanSize),
	}

	go func() {
		ticker := time.NewTicker(time.Millisecond * signalInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if sq.q.Len() > 0 {
					select {
					case sq.C <- struct{}{}:
					default:
					}
				}
			}
		}
	}()

	return sq
}

func (s *SafeQueue) Len() int {
	s.Lock()
	n := s.q.Len()
	s.Unlock()
	return n
}

func (s *SafeQueue) Push(v interface{}) {
	s.Lock()
	defer s.Unlock()
	s.q.Push(v)
}

func (s *SafeQueue) Pop() (interface{}, bool) {
	s.Lock()
	defer s.Unlock()
	return s.q.Pop()
}

func (s *SafeQueue) Front() (interface{}, bool) {
	s.Lock()
	defer s.Unlock()
	return s.q.Front()
}

func main() {
	q := NewSafeQueue()

	var wg sync.WaitGroup

	wg.Add(4)

	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			time.Sleep(time.Second)
			q.Push(i + 1)
		}
	}()

	go func() {
		defer wg.Done()
	LOOP:
		for {
			select {
			case <-q.C:

				for {
					i, ok := q.Pop()
					if !ok {
						continue LOOP
					}
					fmt.Printf("comsumer 1, %d\n", i.(int))
				}

			}
		}
	}()

	go func() {
		defer wg.Done()
	LOOP:
		for {
			select {
			case <-q.C:

				for {
					i, ok := q.Pop()
					if !ok {
						continue LOOP
					}
					fmt.Printf("comsumer 2, %d\n", i.(int))
				}

			}
		}
	}()

	go func() {
		defer wg.Done()
	LOOP:
		for {
			select {
			case <-q.C:

				for {
					i, ok := q.Pop()
					if !ok {
						continue LOOP
					}
					fmt.Printf("comsumer 3, %d\n", i.(int))
				}

			}
		}
	}()

	wg.Wait()

}
