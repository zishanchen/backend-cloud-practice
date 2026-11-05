package main

import "fmt"

func main() {
	// 1. Create a slice
	numbers := []int{10, 20, 30}

	fmt.Println("Initial slice:", numbers)
	fmt.Println("Length of the slice:", len(numbers))
	fmt.Println("Capacity of the slice:", cap(numbers))
	fmt.Printf("Address of the slice: %p\n", &numbers[0])

	fmt.Println()

	// 2. Append within capacity/possibly trigger growth
	numbers = append(numbers, 40)

	fmt.Println("After append 40:", numbers)
	fmt.Println("Length of the slice:", len(numbers))
	fmt.Println("Capacity of the slice:", cap(numbers))
	fmt.Printf("Address of the slice: %p\n", &numbers[0])

	fmt.Println()

	// 3. Keep appending and observe capacity + address
	for i := 50; i <= 100; i += 10 {
		oldCap := cap(numbers)
		oldAddress := &numbers[0]

		numbers = append(numbers, i)

		fmt.Println("After append:", i)
		fmt.Println("slice:", numbers)
		fmt.Println("Length of the slice:", len(numbers))
		fmt.Println("Capacity of the slice:", cap(numbers))
		fmt.Println("Old address:0", oldAddress)
		fmt.Printf("New address: %p\n", &numbers[0])

		if cap(numbers) != oldCap {
			fmt.Println("Capacity changed -> underlying array may have been reallocated")
		}

		fmt.Println()
	}

	// 4. copy()
	source := []int{1, 2, 3}
	destination := make([]int, len(source))

	copied := copy(destination, source)

	fmt.Println("Source:", source)
	fmt.Println("Destination:", destination)
	fmt.Println("Copied elements:", copied)

	// Modify destination
	destination[0] = 99

	fmt.Println("After modifying destination:")
	fmt.Println("Source:", source)
	fmt.Println("Destination:", destination)

}
