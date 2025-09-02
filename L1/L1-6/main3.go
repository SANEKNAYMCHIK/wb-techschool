package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Завершение горутины через контекст
func main() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		counter := 0
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Подсчет окончен")
				return
			default:
				fmt.Println(counter)
				counter++
			}
		}
	}(ctx)
	wg.Wait()
	fmt.Println("Программа завершилась")
}
