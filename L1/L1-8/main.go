package main

import (
	"fmt"
	"math"
)

// Функция, инвертирующая выбранный бит
func invertBit(number int64, bitNum int) int64 {
	bitNum -= 1
	curBit := number >> int64(bitNum)
	if curBit%2 == 0 {
		number = number | (1 << bitNum)
	} else {
		number = number & (math.MaxInt64 ^ (1 << bitNum))
	}
	return number
}

func main() {
	// Изменяемое число
	var number int64 = 217
	// Изменяемый бит
	var bitNum int = 3
	fmt.Printf("Изначальное число в десятичной форме: %d\nВ двоичной форме: %b\n", number, number)
	fmt.Println("=======================")
	number = invertBit(number, bitNum)
	fmt.Printf("Число с измененным %d битом в десятичной форме: %d\nВ двоичной форме: %b\n", bitNum, number, number)

	fmt.Println()

	// Изменяемое число
	var number2 int64 = 62
	// Изменяемый бит
	var bitNum2 int = 4
	fmt.Printf("Изначальное число в десятичной форме: %d\nВ двоичной форме: %b\n", number2, number2)
	fmt.Println("=======================")
	number2 = invertBit(number2, bitNum2)
	fmt.Printf("Число с измененным %d битом в десятичной форме: %d\nВ двоичной форме: %b\n", bitNum2, number2, number2)
}
