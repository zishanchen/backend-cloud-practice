package main

import "fmt"

func calculate(a, b int) (int, int) {
	sum := a + b
	difference := a - b
	return sum, difference
}

func main() {
	sum, difference := calculate(10, 4)
	fmt.Println("Sum:", sum)
	fmt.Println("Difference:", difference)
}
