package main

import (
	"errors"
	"strings"
)

type Contact struct {
	Name  string
	Email string
	Phone string
}

func (c Contact) Summary() string {
	return c.Name + " | " + c.Email + " | " + c.Phone
}

func addContact(
	contacts []Contact,
	name string,
	email string,
	phone string,
) ([]Contact, error) {
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)
	phone = strings.TrimSpace(phone)

	if name == "" {
		return contacts, errors.New("name cannot be empty")
	}

	if email == "" {
		return contacts, errors.New("email cannot be empty")
	}

	if phone == "" {
		return contacts, errors.New("phone cannot be empty")
	}

	contact := Contact{
		Name:  name,
		Email: email,
		Phone: phone,
	}

	contacts = append(contacts, contact)

	return contacts, nil
}

func findContact(contacts []Contact, name string) (Contact, error) {
	name = strings.TrimSpace(name)

	for _, contact := range contacts {
		if strings.EqualFold(contact.Name, name) {
			return contact, nil
		}
	}

	return Contact{}, errors.New("contact not found")
}
