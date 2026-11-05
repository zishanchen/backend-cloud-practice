package main

import (
	"fmt"

	"github.com/zishanchen/backend-cloud-practice/01-go/02-packages-testing/packages/mathutil"
	"github.com/zishanchen/backend-cloud-practice/01-go/02-packages-testing/packages/stringutil"
)

func main() {
	sum := mathutil.Add(5, 3)
	difference := mathutil.Subtract(10, 4)
	product := mathutil.Multiply(6, 7)
	quotient, valid := mathutil.Divide(20, 4)

	if !valid {
		fmt.Println("Error: Division by zero")
		return
	}

	fmt.Println("Sum:", sum)
	fmt.Println("Difference:", difference)
	fmt.Println("Product:", product)
	fmt.Println("Quotient:", quotient)

	text := "hello world"
	upperText := stringutil.ToUpper(text)
	joinedText := stringutil.Join("Hello", "Go")

	fmt.Println("Uppercase:", upperText)
	fmt.Println("Joined Text:", joinedText)
}
