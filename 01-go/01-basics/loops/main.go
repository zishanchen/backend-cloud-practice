package main

import "fmt"

func main() {
	// For loop
	for i := 1; i < 5; i++ {
		fmt.Println("i =", i)
	}

	count := 3

	// For as a while loop
	for count > 0 {
		fmt.Println("Countdown =", count)
		count--
	}

	fmt.Println("Done!")

}
