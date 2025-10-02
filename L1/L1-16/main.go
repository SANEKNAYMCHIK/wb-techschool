package main

import (
	"fmt"
	"math/rand"
)

func quickSort(vals []int) []int {
	res := make([]int, len(vals))
	copy(res, vals)
	quickSortHelper(res, 0, len(res)-1)
	return res
}

func quickSortHelper(vals []int, l, r int) []int {
	if l < r {
		pivotIdx := partition(vals, l, r)
		quickSortHelper(vals, l, pivotIdx-1)
		quickSortHelper(vals, pivotIdx+1, r)
	}
	return vals
}

func partition(nums []int, l, r int) int {
	mid := l + (r-l)/2
	pivot := nums[mid]
	nums[mid], nums[r] = nums[r], nums[mid]
	i := l
	for j := l; j < r; j++ {
		if nums[j] <= pivot {
			nums[i], nums[j] = nums[j], nums[i]
			i++
		}
	}
	nums[i], nums[r] = nums[r], nums[i]
	return i
}

func main() {
	vals := []int{15, 3, 67, 89, 1, 2, 1, 10, 13, 100, 7}
	fmt.Printf("Изначальный массив: %v\n", vals)
	fmt.Printf("Отсортированный массив: %v\n", quickSort(vals))
	fmt.Println("------------------")

	vals2 := []int{20, 19, 17, 15, 8, 7, 5, 3, 2, 1, 0}
	fmt.Printf("Изначальный массив: %v\n", vals2)
	fmt.Printf("Отсортированный массив: %v\n", quickSort(vals2))
	fmt.Println("------------------")

	randomVals := make([]int, 11)
	for i := range randomVals {
		randomVals[i] = rand.Intn(100)
	}
	fmt.Printf("Изначальный массив: %v\n", randomVals)
	fmt.Printf("Отсортированный массив: %v\n", quickSort(randomVals))
}
