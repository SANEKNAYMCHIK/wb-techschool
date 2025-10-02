package main

import (
	"fmt"
	"sync"
)

type syncCounter struct {
	num int
	mu  sync.Mutex
}

func (s *syncCounter) String() string {
	return fmt.Sprintf("%d", s.num)
}

func (s *syncCounter) Add() {
	s.mu.Lock()
	s.num++
	s.mu.Unlock()
}

func main() {
	counter := syncCounter{}
	var wg sync.WaitGroup
	for range 999 {
		wg.Add(1)
		go func(counter *syncCounter) {
			defer wg.Done()
			counter.Add()
		}(&counter)
	}
	wg.Wait()
	fmt.Println(&counter)
}
