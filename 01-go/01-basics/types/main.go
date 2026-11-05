package main

import "fmt"

func main() {
	var age int = 29
	var score float64 = 95.5
	var isStudent bool = true
	var name string = "Zishan"

	fmt.Println("age:", age)
	fmt.Println("score:", score)
	fmt.Println("isStudent:", isStudent)
	fmt.Println("name:", name)

	fmt.Printf("age type: %T\n", age)
	fmt.Printf("score type: %T\n", score)
	fmt.Printf("isStudent type: %T\n", isStudent)
	fmt.Printf("name type: %T\n", name)
}
