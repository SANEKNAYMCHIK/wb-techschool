package main

import "fmt"

func Reverse(val []rune) string {
	n := len(val)
	for i := range n / 2 {
		val[i], val[n-i-1] = val[n-i-1], val[i]
	}
	return string(val)
}

func main() {
	val := "главрыба"
	val2 := "banana"
	val3 := "😁1😂2😃"
	fmt.Println(Reverse([]rune(val)))
	fmt.Println(Reverse([]rune(val2)))
	fmt.Println(Reverse([]rune(val3)))
}
