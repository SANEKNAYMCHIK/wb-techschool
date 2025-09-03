package main

import (
	"fmt"
	"sync"
)

func main() {
	var safeMap sync.Map
	var wg sync.WaitGroup

	for i := 1; i < 10; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			key := fmt.Sprintf("Key%d", val)
			safeMap.Store(key, val)
		}(i)
	}

	wg.Wait()

	for i := range 11 {
		key := fmt.Sprintf("Key%d", i)
		if value, ok := safeMap.Load(fmt.Sprintf("Key%d", i)); ok {
			fmt.Printf("У ключа: %s следующее значение: %d \n", key, value)
		} else {
			fmt.Printf("Ключ %d не найден\n", i)
		}
	}
}
