package main

import "fmt"

func main() {
	vals := []string{"cat", "cat", "dog", "cat", "tree"}
	data := make(map[string]bool)
	for i := range vals {
		data[vals[i]] = true
	}
	res := make([]string, len(data))
	i := 0
	for k := range data {
		res[i] = k
		i++
	}
	fmt.Println(res)
}
