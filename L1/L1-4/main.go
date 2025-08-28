package main

import (
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
			fmt.Println("Введены неккоректные данные для количества воркеров")
			os.Exit(1)
		}
	} else {
		fmt.Print("Введите количество воркеров: ")
		fmt.Scan(&numWorkers)
	}
	jobs := make(chan int, numWorkers*2)
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT)
	var wg sync.WaitGroup
	for range numWorkers {
		wg.Add(1)
		go worker(jobs, &wg)
	}

	go func() {
		counter := 1
		for {
			select {
			case <-stopChan:
				close(jobs)
				return
			default:
				jobs <- counter
				counter++
			}
		}
	}()
	<-stopChan
	// закрываем здесь канал для сигнала, так как канал для получения сигнала - буферизованный на 1 элемент
	// и из-за этого главная горутина читает оттуда сигнал о завершении и переходит к wg.Wait(),
	// но горутина с отправкой не успела прочитать это и продолжает работу и поэтому в главной горутине после чтения сигнала о завершении работы,
	// закрываем канал для сигналов, чтобы горутина отправитель также закрылась и закрыла канал и затем воркеров
	close(stopChan)
	wg.Wait()
}
