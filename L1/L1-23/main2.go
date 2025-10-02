package main

import (
	"fmt"
)

func main() {
	vals := []int{1, 2, 3, 4, 5}
	i := 2
	fmt.Println(vals, len(vals), cap(vals))
	vals = append(vals[:i], vals[i+1:]...)
	fmt.Println(vals, len(vals), cap(vals))
}
