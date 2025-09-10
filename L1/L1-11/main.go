package main

import (
	"fmt"
	"slices"
)

func main() {
	vals1 := []int{8, 14, 1, 2, 17, 2, 2, 6}
	vals2 := []int{2, 5, 143, 12, 9, 14, 3}
	n1 := len(vals1)
	n2 := len(vals2)

	slices.Sort(vals1)
	slices.Sort(vals2)

	var ans []int

	i, j := 0, 0
	for (i < n1) && (j < n2) {
		if vals1[i] > vals2[j] {
			j++
		} else {
			if vals1[i] == vals2[j] {
				if (len(ans) == 0) || ((len(ans) > 0) && (ans[len(ans)-1] != vals1[i])) {
					ans = append(ans, vals1[i])
				}
			}
			i++
		}
	}
	fmt.Println(ans)
}
