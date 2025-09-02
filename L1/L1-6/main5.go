package main

import (
	"fmt"
	"sync"
	"time"
)

// Завершение горутины через time.After
func main() {
	var wg sync.WaitGroup
	timeout := time.After(1 * time.Second)

	wg.Add(1)
	go func(timeout <-chan time.Time) {
		defer wg.Done()
		counter := 0
		for {
			select {
			case <-timeout:
				fmt.Println("Подсчет окончен")
				return
			default:
				fmt.Println(counter)
				counter++
			}
		}
	}(timeout)

	wg.Wait()
	fmt.Println("Программа завершилась")
}
