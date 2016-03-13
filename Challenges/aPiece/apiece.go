package main

import "fmt"

func foo(num ...int) {
	fmt.Println(num)
}
func main() {
	foo(1, 2)

	foo(1, 2, 3)

	aSlice := []int{1, 2, 3, 4}

	foo(aSlice...)

	foo()
}
