package main

import "fmt"

func main() {
	// 1. Declare a fixed-sized array
	var numbers [5]int

	fmt.Println("Empty array:", numbers)

	// 2. Assign values
	numbers[0] = 10
	numbers[1] = 20
	numbers[2] = 30
	numbers[3] = 40
	numbers[4] = 50

	fmt.Println("After assignment:", numbers)

	// 3. Array literal
	scores := [3]int{90, 85, 88}
	fmt.Println("Scores:", scores)

	// 4. Let Go infer the array length
	names := [...]string{"Alice", "Bob", "Charlie"}

	fmt.Println("Names:", names)

	// 5. Iterate with index and value
	for i, v := range scores {
		fmt.Println("Index:", i, "Values:", v)
	}

	// 6. Iterate with name only
	for _, v := range names {
		fmt.Println("Name:", v)
	}

	// 7. Copy an array
	copiedScores := scores

	copiedScores[0] = 100

	fmt.Println("Original scores:", scores)
	fmt.Println("Copied scores:", copiedScores)

	// 8. Compare arrays
	a := [3]int{1, 2, 3}
	b := [3]int{1, 2, 3}

	fmt.Println("a == b", a == b)

	// 9. Get array length
	fmt.Println("Length of numbers:", len(numbers))

}
