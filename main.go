package main

import (
	"fmt"
)

func main() {
	a := make([]int, 0, 2)
	a = add(a)
	fmt.Printf("%p, %v\n", a, a)
	a = append(a, 2)
	fmt.Printf("%p, %v\n", a, a)
}

func add(a []int) []int {
	a = append(a, 1)
	fmt.Printf("%p, %v\n", a, a)
	return a
}
