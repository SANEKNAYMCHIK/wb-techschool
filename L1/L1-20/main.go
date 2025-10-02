package main

import "fmt"

func ReverseWord(vals []rune, i, j int) {
	for i < j {
		vals[i], vals[j] = vals[j], vals[i]
		i++
		j--
	}
}

func ReverseWords(val string) string {
	letters := []rune(val)
	n := len(letters)
	ReverseWord(letters, 0, n-1)
	prevVal := 0
	for i := 0; i <= n; i++ {
		if i == n || letters[i] == ' ' {
			ReverseWord(letters, prevVal, i-1)
			prevVal = i + 1
		}
	}
	return string(letters)
}

func main() {
	val1 := "snow dog sun"
	fmt.Println(ReverseWords(val1))
	val2 := "солнце два три лодка"
	fmt.Println(ReverseWords(val2))
	val3 := "ночь street фонарь аптека"
	fmt.Println(ReverseWords(val3))
}
