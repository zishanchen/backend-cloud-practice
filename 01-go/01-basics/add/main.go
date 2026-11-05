package main

import "fmt"

func add(a, b int) int {
	return a + b
}

func main() {
	var number1 int
	var number2 int

	fmt.Print("Enter first number: ")
	fmt.Scanln(&number1)

	fmt.Print("Enter second number: ")
	fmt.Scanln(&number2)

	result := add(number1, number2)
	fmt.Println("Result:", result)
}
