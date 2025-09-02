package main

import (
	"fmt"
	"sync"
)

// Завершение горутины по условию
func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for counter := range 100 {
			fmt.Println(counter)
		}
		fmt.Println("Подсчет окончен")
	}()
	wg.Wait()
	fmt.Println("Программа завершилась")
}
