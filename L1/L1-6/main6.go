package main

import (
	"fmt"
	"sync"
)

// Завершение горутины через закрытие канала из которого читаем
func main() {
	var wg sync.WaitGroup
	ch := make(chan int)

	wg.Add(1)
	go func(getChan <-chan int) {
		defer wg.Done()
		for {
			val, ok := <-getChan
			if !ok {
				fmt.Println("Подсчет окончен")
				return
			}
			fmt.Println(val)
		}
	}(ch)

	for i := range 115 {
		ch <- i
	}

	close(ch)
	wg.Wait()
	fmt.Println("Программа завершилась")
}
