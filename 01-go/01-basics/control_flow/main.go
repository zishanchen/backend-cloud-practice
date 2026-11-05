package main

import "fmt"

func main() {
	age := 29

	if age < 18 {
		fmt.Println("You are a minor")
	} else if age < 65 {
		fmt.Println("You are an adult")
	} else {
		fmt.Println("You are a senior")
	}

	day := "Sunday"
	switch day {
	case "Saturday", "Sunday":
		fmt.Println("It's the weekend!")
	default:
		fmt.Println("It's a weekday.")
	}
}
