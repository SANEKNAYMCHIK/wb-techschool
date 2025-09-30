package main

import "fmt"

// Интерфейс с пишущей ручкой, который мы ожидаем
type BlackPen interface {
	Write(string)
}

// Существующая ручка с синими чернилами
type BluePen struct{}

func (bp *BluePen) WriteBlue(val string) {
	fmt.Println("*Blue Ink*\n", val)
}

// Адаптер, который работает как BlackPen, но использует BluePen
type PenAdapter struct {
	bluePen *BluePen
}

func (pa *PenAdapter) Write(val string) {
	pa.bluePen.WriteBlue(val)
}

func blackWrite(bp BlackPen) {
	bp.Write("Hello, World!")
}

func main() {
	bluePen := &BluePen{}
	adapter := &PenAdapter{
		bluePen: bluePen,
	}
	blackWrite(adapter)
}
