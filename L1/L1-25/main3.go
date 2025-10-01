package main

import (
	"fmt"
	"time"
)

func sleep(duration int) {
	timeout := time.Now().Add(time.Duration(duration) * time.Second)
	for time.Now().Before(timeout) {
	}
}

func main() {
	startTime := time.Now()
	fmt.Println("Starting timer, current time is:", startTime.Format("15:04:05.000"))
	duration := 5
	sleep(duration)
	fmt.Println("Sleep function ended work, current time is:", time.Now().Format("15:04:05.000"))
	fmt.Println("Duration from program:", duration)
	fmt.Println("Time, that function worked:", time.Since(startTime))
}
