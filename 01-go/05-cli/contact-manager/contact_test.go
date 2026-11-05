package main

import "testing"

func TestAddContact(t *testing.T) {
	tests := []struct {
		name          string
		inputName     string
		inputEmail    string
		inputPhone    string
		expectedLen   int
		expectError   bool
		expectedName  string
		expectedEmail string
		expectedPhone string
	}{
		{
			name:          "valid contact",
			inputName:     "Alice",
			inputEmail:    "alice@example.com",
			inputPhone:    "123456",
			expectedLen:   1,
			expectError:   false,
			expectedName:  "Alice",
			expectedEmail: "alice@example.com",
			expectedPhone: "123456",
		},
		{
			name:        "empty name",
			inputName:   "",
			inputEmail:  "alice@example.com",
			inputPhone:  "123456",
			expectedLen: 0,
			expectError: true,
		},
		{
			name:        "empty email",
			inputName:   "Alice",
			inputEmail:  "",
			inputPhone:  "123456",
			expectedLen: 0,
			expectError: true,
		},
		{
			name:        "empty phone",
			inputName:   "Alice",
			inputEmail:  "alice@example.com",
			inputPhone:  "",
			expectedLen: 0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			contacts := []Contact{}

			result, err := addContact(
				contacts,
				tt.inputName,
				tt.inputEmail,
				tt.inputPhone,
			)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(result) != tt.expectedLen {
				t.Errorf(
					"len(result) = %d; expected %d",
					len(result),
					tt.expectedLen,
				)
			}

			if result[0].Name != tt.expectedName {
				t.Errorf(
					"Name = %s; expected %s",
					result[0].Name,
					tt.expectedName,
				)
			}

			if result[0].Email != tt.expectedEmail {
				t.Errorf(
					"Email = %s; expected %s",
					result[0].Email,
					tt.expectedEmail,
				)
			}

			if result[0].Phone != tt.expectedPhone {
				t.Errorf(
					"Phone = %s; expected %s",
					result[0].Phone,
					tt.expectedPhone,
				)
			}
		})
	}
}

func TestFindContact(t *testing.T) {
	contacts := []Contact{
		{
			Name:  "Alice",
			Email: "alice@example.com",
			Phone: "123456",
		},
		{
			Name:  "Bob",
			Email: "bob@example.com",
			Phone: "987654",
		},
	}

	tests := []struct {
		name         string
		searchName   string
		expectedName string
		expectError  bool
	}{
		{
			name:         "find existing contact",
			searchName:   "Alice",
			expectedName: "Alice",
			expectError:  false,
		},
		{
			name:         "case insensitive search",
			searchName:   "alice",
			expectedName: "Alice",
			expectError:  false,
		},
		{
			name:        "contact not found",
			searchName:  "Charlie",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := findContact(
				contacts,
				tt.searchName,
			)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.Name != tt.expectedName {
				t.Errorf(
					"Name = %s; expected %s",
					result.Name,
					tt.expectedName,
				)
			}
		})
	}
}

func TestContactSummary(t *testing.T) {
	contact := Contact{
		Name:  "Alice",
		Email: "alice@example.com",
		Phone: "123456",
	}

	got := contact.Summary()
	expected := "Alice | alice@example.com | 123456"

	if got != expected {
		t.Errorf(
			"Summary() = %q; expected %q",
			got,
			expected,
		)
	}
}
