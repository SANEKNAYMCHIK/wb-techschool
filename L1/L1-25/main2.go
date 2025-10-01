package main

import (
	"fmt"
	"time"
)

func sleep(duration int) {
	done := make(chan bool)
	go func(done chan bool) {
		time.Sleep(time.Duration(duration) * time.Second)
		done <- true
	}(done)
	<-done
}

func main() {
	startTime := time.Now()
	fmt.Println("Starting timer, current time is:", startTime.Format("15:04:05.000"))
	duration := 3
	sleep(duration)
	fmt.Println("Sleep function ended work, current time is:", time.Now().Format("15:04:05.000"))
	fmt.Println("Duration from program:", duration)
	fmt.Println("Time, that function worked:", time.Since(startTime))
}
