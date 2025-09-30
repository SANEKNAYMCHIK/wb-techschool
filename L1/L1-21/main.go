package main

import "fmt"

// Интерфейс с пишущей ручкой, который мы ожидаем
type Pen interface {
	Write(string)
}

// Структура с кистью
type Brush struct{}

func (bp *Brush) WriteBrush(val string) {
	fmt.Println("*Brush*\n", val)
}

// Адаптер, который работает как Pen, но использует Brush
type PenAdapter struct {
	brush *Brush
}

func (pa *PenAdapter) Write(val string) {
	pa.brush.WriteBrush(val)
}

func WritePen(bp Pen) {
	bp.Write("Hello, World!")
}

func main() {
	brush := &Brush{}
	adapter := &PenAdapter{
		brush: brush,
	}
	WritePen(adapter)
}
