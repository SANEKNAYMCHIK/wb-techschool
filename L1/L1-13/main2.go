package main

import "fmt"

func selfSwap(a, b *int) {
	*a = *a ^ *b
	*b = *a ^ *b
	*a = *a ^ *b
}

func main() {
	a, b := 5, 8
	selfSwap(&a, &b)
	fmt.Println(a, b)
	c, d := 13, -9
	selfSwap(&c, &d)
	fmt.Println(c, d)
	e, f := -19, 2
	selfSwap(&e, &f)
	fmt.Println(e, f)
	g, h := -8, -13
	selfSwap(&g, &h)
	fmt.Println(g, h)
}
