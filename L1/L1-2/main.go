package main

import (
	"fmt"
	"sync"
)

func main() {
	vals := []int{2, 4, 6, 8, 10}
	var wg sync.WaitGroup
	for i := range vals {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			fmt.Println(num * num)
		}(vals[i])
	}
	wg.Wait()
}
