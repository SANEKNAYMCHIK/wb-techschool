package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
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
	const (
		BORDER int = 1000
	)
	if len(os.Args) > 1 {
		numWorkers, err = strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Println("Введены неккоректные данные для количества воркеров")
			os.Exit(1)
		}
	} else {
		fmt.Print("Введите количество воркеров: ")
		fmt.Scan(&numWorkers)
	}
	jobs := make(chan int, numWorkers*2)
	var wg sync.WaitGroup

	for range numWorkers {
		wg.Add(1)
		go worker(jobs, &wg)
	}

	for i := range BORDER {
		jobs <- i
	}
	close(jobs)
	wg.Wait()
}
