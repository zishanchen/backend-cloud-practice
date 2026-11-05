package person

type Person struct {
	Name string
	age  int
}

// Exported function
func NewPerson(name string, age int) Person {
	return Person{
		Name: name,
		age:  age,
	}
}

// Exported method
func (p Person) GetAge() int {
	return p.age
}

// Unexported function
func formatName(name string) string {
	return "Name: " + name
}

// Exported method can use unexported function
func (p Person) Description() string {
	return formatName(p.Name)
}
