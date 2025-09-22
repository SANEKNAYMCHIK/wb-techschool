package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type syncCounter struct {
	num atomic.Int32
}

func (s *syncCounter) String() string {
	return fmt.Sprintf("%d", s.num.Load())
}

func (s *syncCounter) Add() {
	s.num.Add(1)
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
