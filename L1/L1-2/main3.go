package main

import (
	"fmt"
	"runtime"
	"sync/atomic"
)

// Наверное из трех - наихудший вариант работы, так как происходит активное ожидание
// в конце программы в главной горутине и из-за этого остальные горутины могу получать недостаточно времени CPU
// но вызов runtime.Goshed() улучшает ситуацию и позволяет отдавать процессорное время остальным горутинам
func main() {
	vals := []int{2, 4, 6, 8, 10}
	var counter int32

	for i := range vals {
		go func(num int) {
			fmt.Println(num * num)
			atomic.AddInt32(&counter, 1)
		}(vals[i])
	}

	for atomic.LoadInt32(&counter) < int32(len(vals)) {
		runtime.Gosched()
	}
}
