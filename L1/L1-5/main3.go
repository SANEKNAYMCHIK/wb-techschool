package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	ch := make(chan int)
	var wg sync.WaitGroup

	var timeDur time.Duration = 1 // те самые N секунд из задачи
	timeOut := timeDur * time.Second

	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	wg.Add(1)
	go func(ctx context.Context, setChan chan<- int) {
		defer wg.Done()
		counter := 1
		for {
			select {
			case setChan <- counter:
				counter++
			case <-ctx.Done():
				close(setChan)
				return
			}
		}
	}(ctx, ch)

	wg.Add(1)
	go func(ctx context.Context, getChan <-chan int) {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				fmt.Printf("Истекло %d секунд\n", timeDur)
				return
			case val := <-getChan:
				fmt.Println(val)
			}
		}
	}(ctx, ch)
	<-ctx.Done()
	fmt.Println("Завершение программы")
	wg.Wait()
}
