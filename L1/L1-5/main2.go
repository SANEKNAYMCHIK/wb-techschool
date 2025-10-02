package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	ch := make(chan int)
	done := make(chan struct{})
	var wg sync.WaitGroup

	var timeDur time.Duration = 1 // те самые N секунд из задачи
	timeOut := timeDur * time.Second
	timeout := time.After(timeOut)

	wg.Add(1)
	go func(setChan chan<- int, done chan struct{}) {
		defer wg.Done()
		counter := 1
		for {
			select {
			case setChan <- counter:
				counter++
			case <-done:
				close(setChan)
				return
			}
		}
	}(ch, done)

	wg.Add(1)
	go func(getChan <-chan int, done chan struct{}) {
		defer wg.Done()
		for {
			select {
			case <-done:
				fmt.Printf("Истекло %d секунд\n", timeDur)
				return
			case val := <-getChan:
				fmt.Println(val)
			}
		}
	}(ch, done)

	<-timeout
	close(done)
	fmt.Println("Завершение программы")
	wg.Wait()
}
