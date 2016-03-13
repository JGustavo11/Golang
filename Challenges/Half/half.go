package main

import "fmt"

func half (x int) (int, bool) {
	return x/2, x%2==0
}

func main() {
	var z int
	fmt.Print("Enter a number ")
	fmt.Scan(&z)
	y, w := half(z)
	fmt.Println(y, w)

}
