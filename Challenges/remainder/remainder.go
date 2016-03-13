package main

import "fmt"

func main() {
	var num1 int
	var num2 int


	fmt.Print("Enter a number: ")
	fmt.Scan(&num1)

	fmt.Print("Enter a another number, preferably smaller ")
	fmt.Scan(&num2)


	fmt.Println("The remainder is: ", num1%num2)
}
