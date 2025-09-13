package main

import (
	"fmt"
	"reflect"
)

func getType(val interface{}) {
	switch v := val.(type) {
	case int:
		fmt.Printf("value: %d | type: int\n", v)
	case string:
		fmt.Printf("value: %s | type: string\n", v)
	case bool:
		fmt.Printf("value: %t | type: bool\n", v)
	default:
		t := reflect.TypeOf(val)
		if t != nil && t.Kind() == reflect.Chan {
			fmt.Printf("value: %v | type: %v\n", v, t)
		}
	}
}

func main() {
	getType(13)
	getType("hi")
	getType(true)
	ch := make(chan int)
	getType(ch)
}
