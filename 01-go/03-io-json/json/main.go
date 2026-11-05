package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Contact struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

func encodeJSON(contact Contact) ([]byte, error) {
	data, err := json.MarshalIndent(contact, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("encode JSON: %w", err)
	}

	return data, nil
}

func decodeJSON(data []byte) (Contact, error) {
	var contact Contact

	err := json.Unmarshal(data, &contact)
	if err != nil {
		return Contact{}, fmt.Errorf("decode JSON: %w", err)
	}

	return contact, nil
}

func saveJSON(filename string, data []byte) error {
	err := os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("save JSON file: %w", err)
	}

	return nil
}

func loadJSON(filename string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("load JSON file: %w", err)
	}

	return data, nil
}

func main() {
	contact := Contact{
		Name:  "Zishan",
		Email: "zishan@example.com",
		Phone: "123456789",
	}

	// 1. Encode struct -> JSON
	data, err := encodeJSON(contact)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Encoded JSON:")
	fmt.Println(string(data))

	// 2. Save JSON to file
	filename := "contact.json"

	err = saveJSON(filename, data)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("\nJSON saved to", filename)

	// 3. Read JSON from file
	fileData, err := loadJSON(filename)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// 4. Decode JSON -> struct
	decodedContact, err := decodeJSON(fileData)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("\nDecoded Contact:")
	fmt.Println("Name:", decodedContact.Name)
	fmt.Println("Email:", decodedContact.Email)
	fmt.Println("Phone:", decodedContact.Phone)
}
