package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
)

func worker(jobs <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range jobs {
		fmt.Println(job)
	}
}

func main() {
	var numWorkers int
	var err error
	if len(os.Args) > 1 {
		numWorkers, err = strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Println("Введены некорректные данные для количества воркеров")
			os.Exit(1)
		}
	} else {
		fmt.Print("Введите количество воркеров: ")
		fmt.Scan(&numWorkers)
	}

	// создаем контекст, который автоматически отменяется при получении ctrl+c
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer stop()

	jobs := make(chan int, numWorkers*2)
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(jobs, &wg)
	}

	go func() {
		counter := 1
		for {
			select {
			case <-ctx.Done():
				close(jobs)
				return
			default:
				jobs <- counter
				counter++
			}
		}
	}()
	<-ctx.Done()
	wg.Wait()
}
