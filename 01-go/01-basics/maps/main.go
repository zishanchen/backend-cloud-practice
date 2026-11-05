package main

import "fmt"

func main() {
	// 1. Create a map
	ages := map[string]int{
		"Alice": 25,
		"Bob":   30,
	}

	fmt.Println("Initial map:", ages)

	// 2. Create / Update
	ages["Charlie"] = 28
	ages["Alice"] = 26

	fmt.Println("After create/update:", ages)

	// 3. Read
	fmt.Println("Bob's age:", ages["Bob"])

	// 4. comma-ok: check whether a key exists
	age, ok := ages["David"]

	if ok {
		fmt.Println("David's age:", age)
	} else {
		fmt.Println("David does not exist")
	}

	// 5. Delete
	delete(ages, "Bob")

	fmt.Println("After deleting Bob's age:", ages)

	// 6. Iterate over a map
	fmt.Println("All entries:")

	for name, age := range ages {
		fmt.Println(name, "->", age)
	}

	// 7. Frequency counter
	words := []string{
		"go",
		"java",
		"go",
		"python",
		"go",
		"java",
	}

	frequency := make(map[string]int)

	for _, word := range words {
		frequency[word]++
	}

	fmt.Println("Word Frequency:")

	for word, count := range frequency {
		fmt.Println(word, "appears", count, "times")
	}

}
