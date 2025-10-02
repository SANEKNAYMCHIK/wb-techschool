package main

import (
	"fmt"
	"math/big"
)

func main() {
	string1 := "593475236412374345634856345643583457358832490"
	string2 := "123456431235235623432348723573298532985732924"
	val1, err := new(big.Int).SetString(string1, 10)
	if !err {
		fmt.Println("Error! Something went wrong with input value")
	}
	val2, err := new(big.Int).SetString(string2, 10)
	if !err {
		fmt.Println("Error! Something went wrong with input value")
	}

	sumRes := new(big.Int).Add(val1, val2)
	fmt.Printf("Sum %d + %d = %d\n", val1, val2, sumRes)

	diffRes := new(big.Int).Sub(val1, val2)
	fmt.Printf("Difference %d - %d = %d\n", val1, val2, diffRes)

	prodRes := new(big.Int).Mul(val1, val2)
	fmt.Printf("Multiplication %d * %d = %d\n", val1, val2, prodRes)

	if val2.Sign() == 0 {
		fmt.Println("Division by zero!")
	} else {
		divRes := new(big.Int).Div(val1, val2)
		fmt.Printf("Division %d / %d = %d\n", val1, val2, divRes)
	}
}
