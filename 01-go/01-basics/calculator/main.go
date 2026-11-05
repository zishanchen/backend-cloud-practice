package main

import "fmt"

func calculate(a, b float64, operator string) (float64, bool) {
	switch operator {
	case "+":
		return a + b, true
	case "-":
		return a - b, true
	case "*":
		return a * b, true
	case "/":
		if b == 0 {
			return 0, false
		}
		return a / b, true
	default:
		return 0, false
	}
}

func getUserInputs() (bool, float64, float64, string) {
	var a, b float64
	var operator string

	fmt.Print("Enter first number: ")
	_, err := fmt.Scan(&a)

	if err != nil {
		fmt.Println("Invalid first number")
		return false, 0, 0, ""
	}

	fmt.Print("Enter operator (+, -, *, /): ")
	_, err = fmt.Scan(&operator)

	if err != nil {
		fmt.Println("Invalid operator")
		return false, 0, 0, ""
	}

	fmt.Print("Enter second number: ")
	_, err = fmt.Scan(&b)

	if err != nil {
		fmt.Println("Invalid second number")
		return false, 0, 0, ""
	}

	return true, a, b, operator
}

func main() {
	ok, a, b, operator := getUserInputs()
	if !ok {
		fmt.Println("Invalid input. Exiting.")
		return
	}
	result, valid := calculate(float64(a), float64(b), operator)
	if !valid {
		fmt.Println("Invalid operation or division by zero.")
		return
	}

	fmt.Printf("Result: %.2f\n", result)
}
