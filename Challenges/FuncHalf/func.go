package main

import "fmt"

func main() {
	half := func (x int) (int, bool){
		return x/2, x%2==0
	}
	var z int
	fmt.Print("Enter a number: ")
	fmt.Scan(&z)
	fmt.Println(half(z))

}
