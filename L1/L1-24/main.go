package main

import (
	"fmt"
	"math"
)

type Point struct {
	x float64
	y float64
}

func NewPoint(x, y float64) *Point {
	return &Point{
		x: x,
		y: y,
	}
}

func (p *Point) Distance(other *Point) float64 {
	return math.Sqrt(math.Pow(p.x-other.x, 2) + math.Pow(p.y-other.y, 2))
}

func main() {
	point1 := NewPoint(0, 0)
	point2 := NewPoint(3, 4)
	fmt.Printf("Distance between point1 and point2 is %.3f\n", point1.Distance(point2))

	point3 := NewPoint(2, 2)
	point4 := NewPoint(-2, -2)
	fmt.Printf("Distance between point3 and point4 is %.3f\n", point3.Distance(point4))

	point5 := NewPoint(1, 7)
	point6 := NewPoint(10, 15)
	fmt.Printf("Distance between point5 and point6 is %.3f\n", point5.Distance(point6))
}
