package main

import (
	"fmt"
	"unicode"
)

func UniqueElems(val string) bool {
	data := make(map[rune]bool)
	for _, elem := range val {
		lower := unicode.ToLower(elem)
		if data[lower] {
			return false
		}
		data[lower] = true
	}
	return true
}

func main() {
	fmt.Println(UniqueElems("abcd"))
	fmt.Println(UniqueElems("abCdefAaf"))
	fmt.Println(UniqueElems("aabcd"))
	fmt.Println(UniqueElems("abdfjreoiA"))
	fmt.Println(UniqueElems("b"))
	fmt.Println(UniqueElems("Bb"))
	fmt.Println(UniqueElems("VfwhF"))
}
