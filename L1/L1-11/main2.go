package main

import "fmt"

func main() {
	vals1 := []int{8, 14, 1, 2, 17, 2, 2, 6, 3, 3}
	vals2 := []int{2, 5, 5, 143, 12, 9, 14, 3}
	data := make(map[int]bool)

	ans := make([]int, 0)

	for i := range vals1 {
		if _, ok := data[vals1[i]]; !ok {
			data[vals1[i]] = false
		}
	}
	for i := range vals2 {
		if val, ok := data[vals2[i]]; ok {
			if !val {
				ans = append(ans, vals2[i])
			}
			data[vals2[i]] = true
		}
	}
	fmt.Println(ans)
}
