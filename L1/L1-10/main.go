package main

import "fmt"

func main() {
	vals := []float32{-25.4, -27.0, 13.0, 19.0, 15.5, 24.5, -21.0, 32.5, 8, 0.1, 10}
	data := make(map[int][]float32)
	for i := range vals {
		diap := int(vals[i]) / 10 * 10
		data[diap] = append(data[diap], vals[i])
	}
	for k, v := range data {
		fmt.Printf("key %d : %v\n", k, v)
	}
}
