# Contact Manager CLI

A simple command-line contact manager written in Go.

This project is part of my Go learning practice and focuses on separating CLI input/output from business logic and writing testable code.

## Features

- Add a contact
- List all contacts
- Find a contact by name
- Case-insensitive contact search
- Input validation
- Error handling
- Table-driven unit tests

## Project Structure

```text
contact-manager/
├── main.go
├── contact.go
├── contact_test.go
└── README.md
```

### `main.go`

Handles:

- CLI menu
- User input
- Terminal output

### `contact.go`

Contains the core business logic:

- `Contact` struct
- `Summary()`
- `addContact()`
- `findContact()`

### `contact_test.go`

Contains unit tests for:

- Adding valid contacts
- Invalid input
- Whitespace handling
- Contact search
- Case-insensitive search
- Contact summary

## Run

From the repository root:

```bash
go run ./cli/contact-manager
```

Example:

```text
=== Contact Manager ===
1. Add Contact
2. List Contacts
3. Find Contact
4. Exit
Choose an option:
```

## Tests

Run all tests:

```bash
go test ./cli/contact-manager
```

Run tests with detailed output:

```bash
go test -v ./cli/contact-manager
```

## Concepts Practiced

- Structs
- Methods
- Slices
- Functions
- Error handling
- `bufio.Reader`
- String handling
- Separation of concerns
- Table-driven tests
- Subtests with `t.Run()`

## Design

The CLI layer and business logic are separated.

```text
CLI input/output
      ↓
main.go
      ↓
business logic
      ↓
contact.go
      ↓
unit tests
```

This makes the core logic easier to test and reuse.

## Current Limitations

Contacts are currently stored only in memory.

They are lost when the program exits.

A future version can add JSON persistence and additional CRUD operations.