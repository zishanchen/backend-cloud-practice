package main

import "fmt"

type Person struct {
	Name string
	Age  int
}

type Contact struct {
	Email string
	Phone string
}

// Value receiver:
// only reads data, does not modify the original struct

func (p Person) Introduce() string {
	return fmt.Sprintf("Hi, I am %s and I am %d years old", p.Name, p.Age)
}

// Pointer receiver:
// modifies the original struct
func (p *Person) HaveBirthday() {
	p.Age++
}

// Pointer receiver:
// updates the person's name
func (p *Person) Rename(newName string) {
	p.Name = newName
}

// Value receiver
// returns formatted contact information
func (c Contact) Summary() string {
	return fmt.Sprintf("Email: %s, Phone: %s", c.Email, c.Phone)
}

// Pointer Receiver
// modifies the original contact
func (c Contact) UpdateEmail(newEmail string) {
	c.Email = newEmail
}

func (c Contact) UpdatePhone(newPhone string) {
	c.Phone = newPhone
}

func main() {
	person := Person{
		Name: "Zishan",
		Age:  29,
	}

	contact := Contact{
		Email: "zishanc@example.com",
		Phone: "123456789",
	}

	// Value receiver
	fmt.Println(person.Introduce())

	// Pointer receiver
	person.HaveBirthday()
	fmt.Println("After birthday:", person.Age)

	person.Rename("Zishan Chen")
	fmt.Println("After rename:", person.Name)

	// Contact methods
	fmt.Println(contact.Summary())

	contact.UpdateEmail("zishan.chen@example.com")
	contact.UpdatePhone("987654321")

	fmt.Println("Updated contact:")
	fmt.Println(contact.Summary())
}
