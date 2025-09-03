package main

import (
	"fmt"
	"sync"
)

type SafeMap[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]V
}

func NewSafeMap[K comparable, V any]() *SafeMap[K, V] {
	return &SafeMap[K, V]{
		data: make(map[K]V),
	}
}

func (s *SafeMap[K, V]) Set(key K, value V) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

func (s *SafeMap[K, V]) Get(key K) (V, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok := s.data[key]
	return value, ok
}

func main() {
	safeMap := NewSafeMap[string, int]()
	var wg sync.WaitGroup

	for i := 1; i < 10; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			key := fmt.Sprintf("Key%d", val)
			safeMap.Set(key, val)
		}(i)
	}

	wg.Wait()

	for i := range 11 {
		key := fmt.Sprintf("Key%d", i)
		if value, ok := safeMap.Get(fmt.Sprintf("Key%d", i)); ok {
			fmt.Printf("У ключа: %s следующее значение: %d \n", key, value)
		} else {
			fmt.Printf("Ключ %d не найден\n", i)
		}
	}
}
