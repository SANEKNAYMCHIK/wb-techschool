package main

import "fmt"

func binSearch(nums []int, elem int) int {
	l, r := 0, len(nums)-1
	for l <= r {
		mid := l + (r-l)/2
		if nums[mid] < elem {
			l = mid + 1
		} else if nums[mid] > elem {
			r = mid - 1
		} else {
			return mid
		}
	}
	return -1
}

func main() {
	vals := []int{1, 3, 5, 6, 10, 15, 18, 21}
	fmt.Println(binSearch(vals, 5))
	fmt.Println(binSearch(vals, 21))
	fmt.Println(binSearch(vals, 19))
}
