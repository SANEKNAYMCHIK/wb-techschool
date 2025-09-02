package main

import (
	"fmt"
	"sync"
	"time"
)

// Завершение горутины через канал уведомлений
func main() {
	var wg sync.WaitGroup
	done := make(chan struct{})

	wg.Add(1)
	go func(done <-chan struct{}) {
		defer wg.Done()
		counter := 0
		for {
			select {
			case <-done:
				fmt.Println("Подсчет окончен")
				return
			default:
				fmt.Println(counter)
				counter++
			}
		}
	}(done)
	time.Sleep(1 * time.Second)
	close(done)
	wg.Wait()
	fmt.Println("Программа завершилась")
}
