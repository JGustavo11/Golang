package main

import "fmt"

func great(numbers ...int) int {
	var large int

	for _, x := range numbers {

		if x > large {

			large = x
		}
	}
	return large
}

func main() {
	max := great(150, 110, 45, 175, 23, 12, 77)
	fmt.Println(max)
}
