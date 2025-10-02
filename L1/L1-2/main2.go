package main

import "fmt"

func main() {
	vals := []int{2, 4, 6, 8, 10}
	done := make(chan bool, len(vals))
	for i := range vals {
		go func(num int, ch chan<- bool) {
			fmt.Println(num * num)
			ch <- true
		}(vals[i], done)
	}
	for i := 0; i < len(vals); i++ {
		<-done
	}
}
