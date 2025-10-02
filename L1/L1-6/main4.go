package main

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

// Завершение горутины через runtime.Goexit()
func main() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		counter := 0
		for {
			fmt.Println(counter)
			counter++
			if counter > 115 {
				fmt.Println("Подсчет окончен")
				runtime.Goexit()
			}
		}
	}(ctx)

	wg.Wait()
	fmt.Println("Программа завершилась")
}
