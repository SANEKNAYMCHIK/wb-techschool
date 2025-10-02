package main

import (
	"fmt"
)

func getType(val interface{}) {
	switch v := val.(type) {
	case int:
		fmt.Printf("value: %d | type: int\n", v)
	case string:
		fmt.Printf("value: %s | type: string\n", v)
	case bool:
		fmt.Printf("value: %t | type: bool\n", v)
	case chan int:
		fmt.Printf("value: %v | type: chan int\n", v)
	case chan string:
		fmt.Printf("value: %v | type: chan string\n", v)
	case chan bool:
		fmt.Printf("value: %v | type: chan bool\n", v)
	}
}

func main() {
	getType(13)
	getType("hi")
	getType(true)
	ch := make(chan int)
	getType(ch)
}
