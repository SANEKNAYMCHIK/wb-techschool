package main

import (
	"fmt"
	"sync"
)

func getVal(chanPush <-chan int, chanPull chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for item := range chanPush {
		chanPull <- item * 2
	}
}

func printVal(chanPull <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for item := range chanPull {
		fmt.Println("Получили новое значение x*2=", item)
	}
}

func main() {
	chanPush := make(chan int)
	chanPull := make(chan int)

	var wgWorkers sync.WaitGroup
	var wgPrinter sync.WaitGroup

	vals := []int{1, 5, 7, 13, 19}
	amountWorkers := 3

	for range amountWorkers {
		wgWorkers.Add(1)
		go getVal(chanPush, chanPull, &wgWorkers)
	}
	wgPrinter.Add(1)
	go printVal(chanPull, &wgPrinter)

	go func(chanPush chan<- int) {
		for _, elem := range vals {
			chanPush <- elem
		}
		close(chanPush)
	}(chanPush)
	wgWorkers.Wait()
	close(chanPull)
	wgPrinter.Wait()
	fmt.Println("Программа завершена")
}
