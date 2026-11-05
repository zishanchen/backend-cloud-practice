package main

import "fmt"

// Define an interface
type Greeter interface {
	Greet() string
}

type Person struct {
	Name string
}

type Contact struct {
	Email string
}

// Person implements Greeter
func (p Person) Greet() string {
	return "Hello, I'm " + p.Name
}

// Contact also implements Greeter
func (c Contact) Greet() string {
	return "Contact email: " + c.Email
}

// Function accepts any value that implements Greeter
func printGreeting(g Greeter) {
	fmt.Println(g.Greet())
}

func main() {
	person := Person{
		Name: "Zishan",
	}

	contact := Contact{
		Email: "zishan@example.com",
	}

	printGreeting(person)
	printGreeting(contact)

	// Interface variable
	var greeter Greeter

	greeter = person
	fmt.Println(greeter.Greet())

	greeter = contact
	fmt.Println(greeter.Greet())
}
