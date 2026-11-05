package main

import (
	"fmt"

	"github.com/zishanchen/backend-cloud-practice/01-go/02-packages-testing/visibility/person"
)

func main() {
	p := person.NewPerson("Zishan", 29)

	// Exported field
	fmt.Println("Name:", p.Name)

	// Unexported field cannot be accessed directly
	// fmt.Println(p.age) // ERROR

	// Access unexported data through exported method
	fmt.Println("Age:", p.GetAge())

	// Exported method
	fmt.Println(p.Description())

	// Unexported function cannot be called from another package
	// person.formatName("Zishan") // ERROR
}
