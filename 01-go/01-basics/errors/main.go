package main

import (
	"errors"
	"fmt"
)

// errors.New: fixed error message
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("Division by zero")
	}

	return a / b, nil
}

// fmt.Errorf: dynamic error message
func withdraw(balance, amount float64) (float64, error) {
	if amount <= 0 {
		return balance, fmt.Errorf("Invalid withdrawl amount: %.2f", amount)
	}

	if amount > balance {
		return balance, fmt.Errorf("Insufficient balance: balance=%.2f, requested=%.2f", balance, amount)
	}

	return balance - amount, nil
}

func main() {
	// 1. Successful call
	result, err := divide(10, 2)

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("10 / 2:", result)
	}

	// 2. errors.New example
	result, err = divide(10, 0)

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("10 / 0:", result)
	}

	// 3. fmt.Errorf example
	balance := 100.0

	newBalance, err := withdraw(balance, 30)

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("New balance:", newBalance)
	}

	// 4. Dynamic error message
	newBalance, err = withdraw(balance, 150)

	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("New balance:", newBalance)
	}
}
