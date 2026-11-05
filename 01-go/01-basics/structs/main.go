package main

import "fmt"

type Address struct {
	City    string
	Country string
}

type Person struct {
	Name    string
	Age     int
	Student bool
	Address Address
}

func main() {
	// 1. Initialize with field names
	person1 := Person{
		Name:    "Zishan",
		Age:     29,
		Student: true,
		Address: Address{
			City:    "Munich",
			Country: "Germany",
		},
	}

	fmt.Println("Person 1:", person1)

	// 2. Access fields
	fmt.Println("Name:", person1.Name)
	fmt.Println("City:", person1.Address.City)

	// 3. Update fields
	person1.Age = 30
	person1.Address.City = "Berlin"

	fmt.Println("Updated age:", person1.Age)
	fmt.Println("Updated city:", person1.Address.City)

	// 4. Initialize an empty struct
	var person2 Person

	person2.Name = "Alice"
	person2.Age = 25
	person2.Student = false

	fmt.Println("Person 2:", person2)

	// 5. Another initialization style
	person3 := Person{
		Name: "Bob",
		Age:  28,
	}

	fmt.Println("Person 3:", person3)
}
