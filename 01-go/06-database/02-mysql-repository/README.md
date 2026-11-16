# MySQL Repository Practice

This exercise builds a small Go HTTP API with a layered architecture:

```text
HTTP
↓
Handler
↓
Service
↓
Repository
↓
database/sql
↓
MySQL
```

The main goals are:

- separate HTTP logic from database access
- introduce a Service layer for business logic
- implement Repository methods for `Create`, `Get`, and `List`
- use parameterized MySQL queries with `?`
- avoid SQL injection
- use `context.Context` with database operations
- handle `sql.ErrNoRows`
- write table-driven tests for Handler and Service layers
- test without requiring a real MySQL database

---

## Project Structure

```text
02-mysql-repository/
├── main.go
├── model.go
├── handler.go
├── handler_test.go
├── service.go
├── service_test.go
├── repository.go
└── README.md
```

The SQL schema is stored in:

```text
02-sql/mysql-task-api-schema.sql
```

---

# Architecture

```text
Client
↓
HTTP Request
↓
TaskHandler
↓
TaskService
↓
TaskRepository
↓
MySQLTaskRepository
↓
database/sql
↓
MySQL
```

Each layer has a different responsibility.

---

# Handler Layer

The Handler is responsible for HTTP concerns.

Examples:

```text
decode JSON
validate HTTP input
parse path parameters
map errors to HTTP status codes
encode JSON responses
```

The Handler should call:

```go
h.service.Create(...)
h.service.Get(...)
h.service.List(...)
```

It should not contain SQL such as:

```go
db.Exec(...)
db.Query(...)
db.QueryRow(...)
```

The Handler should not know how MySQL stores tasks.

---

# Service Layer

The Service contains application and business logic.

Example:

```go
func (s *taskService) Create(
	ctx context.Context,
	input CreateTaskRequest,
) (Task, error) {
	input.Title = strings.TrimSpace(input.Title)

	if input.Title == "" {
		return Task{}, ErrTitleRequired
	}

	return s.repo.Create(
		ctx,
		input,
	)
}
```

The Service handles rules such as:

```text
trim task title
validate required title
```

The Handler does not need to know how those rules are implemented.

---

# Repository Layer

The Repository is responsible for database access.

The interface:

```go
type TaskRepository interface {
	Create(
		ctx context.Context,
		input CreateTaskRequest,
	) (Task, error)

	Get(
		ctx context.Context,
		id int64,
	) (Task, error)

	List(
		ctx context.Context,
	) ([]Task, error)
}
```

The MySQL implementation uses:

```text
database/sql
+
go-sql-driver/mysql
```

---

# Dependency Flow

At application startup:

```go
db, err := openDB()
```

Then:

```go
repo := NewMySQLTaskRepository(db)
```

Then:

```go
service := NewTaskService(repo)
```

Then:

```go
handler := NewTaskHandler(service)
```

So the dependency chain is:

```text
sql.DB
↓
MySQLTaskRepository
↓
TaskService
↓
TaskHandler
↓
ServeMux
```

---

# MySQL Schema

The API uses:

```sql
CREATE TABLE IF NOT EXISTS api_tasks (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(255) NOT NULL,
    done BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

Fields:

```text
id
→ primary key

title
→ task title

done
→ completion state

created_at
→ creation timestamp
```

---

# Repository Create

The Create method uses:

```sql
INSERT INTO api_tasks (title, done)
VALUES (?, ?)
```

Go executes it with:

```go
result, err := r.db.ExecContext(
	ctx,
	query,
	input.Title,
	input.Done,
)
```

Then:

```go
id, err := result.LastInsertId()
```

returns the generated task ID.

Flow:

```text
CreateTaskRequest
↓
Repository.Create
↓
ExecContext
↓
INSERT
↓
LastInsertId
↓
Task
```

---

# Repository Get

The Get method uses:

```sql
SELECT id, title, done
FROM api_tasks
WHERE id = ?
```

Go executes:

```go
err := r.db.QueryRowContext(
	ctx,
	query,
	id,
).Scan(
	&task.ID,
	&task.Title,
	&task.Done,
)
```

`QueryRowContext` is appropriate when one row is expected.

---

# `sql.ErrNoRows`

If no task exists:

```go
Scan(...)
```

returns:

```go
sql.ErrNoRows
```

The Repository translates this into:

```go
ErrTaskNotFound
```

Example:

```go
if errors.Is(err, sql.ErrNoRows) {
	return Task{}, ErrTaskNotFound
}
```

The Handler then maps:

```text
ErrTaskNotFound
↓
404 Not Found
```

So:

```text
database-specific error
↓
application error
↓
HTTP response
```

---

# Repository List

The List method uses:

```sql
SELECT id, title, done
FROM api_tasks
ORDER BY id
```

Go executes:

```go
rows, err := r.db.QueryContext(
	ctx,
	query,
)
```

Then:

```go
defer rows.Close()

for rows.Next() {
	var task Task

	if err := rows.Scan(
		&task.ID,
		&task.Title,
		&task.Done,
	); err != nil {
		...
	}

	tasks = append(
		tasks,
		task,
	)
}
```

After iteration:

```go
if err := rows.Err(); err != nil {
	...
}
```

The typical multi-row pattern is:

```text
QueryContext
↓
rows.Next
↓
rows.Scan
↓
append
↓
rows.Err
```

---

# SQL Injection Prevention

Do not build SQL by concatenating user input.

Bad:

```go
query :=
	"SELECT * FROM api_tasks WHERE id = " +
	userInput
```

Bad:

```go
query := fmt.Sprintf(
	"SELECT * FROM api_tasks WHERE id = %s",
	userInput,
)
```

Instead, use placeholders:

```sql
SELECT id, title, done
FROM api_tasks
WHERE id = ?
```

and pass values separately:

```go
r.db.QueryRowContext(
	ctx,
	query,
	id,
)
```

For MySQL:

```text
?
```

is the parameter placeholder.

---

# Why Parameterized Queries Are Safer

Suppose the user input is:

```text
x' OR 1=1 --
```

With:

```sql
INSERT INTO api_tasks (title, done)
VALUES (?, ?)
```

the database receives:

```text
SQL structure
+
separate parameter values
```

The input remains data rather than becoming part of the SQL syntax.

So:

```text
x' OR 1=1 --
```

is stored as a normal title.

---

# Context

Repository methods accept:

```go
ctx context.Context
```

The Handler passes:

```go
r.Context()
```

Flow:

```text
HTTP request
↓
request context
↓
Handler
↓
Service
↓
Repository
↓
ExecContext / QueryContext
```

This allows database operations to follow request:

```text
cancellation
timeout
lifecycle
```

---

# Error Wrapping

Repository errors are wrapped with additional context:

```go
return Task{}, fmt.Errorf(
	"create task: %w",
	err,
)
```

`%w` preserves the original error.

This allows:

```go
errors.Is(...)
```

to keep working while adding useful context.

---

# Client Errors vs Server Errors

The HTTP client should receive a safe message:

```json
{
  "error": "failed to create task"
}
```

The server can log the detailed error:

```go
log.Printf(
	"create task failed: %v",
	err,
)
```

Pattern:

```text
Client
→ safe error message

Server log
→ detailed internal error
```

---

# Testing Strategy

The Handler and Service layers are tested without a real MySQL database.

Testing architecture:

```text
Handler Test
↓
Fake TaskService
```

and:

```text
Service Test
↓
Fake TaskRepository
```

This avoids:

```text
starting MySQL
creating test tables
network/database dependencies
```

for unit tests.

---

# Table-Driven Tests

The tests use the standard table-driven pattern:

```go
tests := []struct {
	name string
	// input
	// expected result
}{
	{
		name: "case one",
	},
	{
		name: "case two",
	},
}

for _, tc := range tests {
	t.Run(
		tc.name,
		func(t *testing.T) {
			// arrange
			// act
			// assert
		},
	)
}
```

This allows the same test logic to run against many scenarios.

---

# Handler Tests

Handler tests use:

```go
httptest.NewRequest(...)
```

to create an HTTP request.

And:

```go
httptest.NewRecorder()
```

to capture the HTTP response.

The Handler tests cover scenarios such as:

```text
successful task creation
invalid JSON
empty-title validation
service failure

successful task lookup
invalid task ID
task not found
service failure

successful task listing
service failure
```

Example flow:

```text
fake HTTP request
↓
TaskHandler
↓
fake TaskService
↓
httptest.ResponseRecorder
↓
assert HTTP status/body
```

---

# Service Tests

Service tests use a fake Repository.

Covered scenarios include:

```text
successful creation
title trimming
empty-title validation
repository error

successful lookup
task not found

successful list
repository list error
```

Example:

```text
CreateTaskRequest
↓
TaskService
↓
Fake Repository
↓
assert result/error
```

---

# Why Use Fakes?

Without fakes:

```text
Handler test
↓
Service
↓
Repository
↓
MySQL
```

A Handler test would depend on the entire stack.

With fakes:

```text
Handler test
↓
Fake Service
```

This isolates the Handler.

Similarly:

```text
Service test
↓
Fake Repository
```

isolates the Service.

The rule is:

> Test one layer without requiring all lower infrastructure layers.

---

# Run Tests

Run this package:

```bash
go test -v \
  ./01-go/06-database/02-mysql-repository
```

Run all Go tests in the repository:

```bash
go test ./...
```

Expected result:

```text
no FAIL
```

Packages without tests may show:

```text
[no test files]
```

which is normal.

---

# Coverage

Check statement coverage:

```bash
go test \
  -cover \
  ./01-go/06-database/02-mysql-repository
```

Generate a coverage profile:

```bash
go test \
  -coverprofile=coverage.out \
  ./01-go/06-database/02-mysql-repository
```

Inspect coverage by function:

```bash
go tool cover \
  -func=coverage.out
```

---

# Current Coverage Scope

The current unit tests focus on:

```text
Handler layer
Service layer
```

Specifically:

```text
HTTP success paths
HTTP validation errors
HTTP status mapping
Service business validation
Service error propagation
Repository/service dependency calls
```

The unit tests intentionally do not test:

```text
real MySQL queries
database connectivity
MySQL driver behavior
Repository integration against a live database
main.go startup wiring
```

Those belong to:

```text
integration tests
```

rather than these Handler/Service unit tests.

---

# Unit Test vs Integration Test

Current tests:

```text
Handler + fake Service
Service + fake Repository
```

are unit tests.

A future test such as:

```text
MySQLTaskRepository
↓
real MySQL
```

would be an integration test.

Difference:

```text
Unit Test
→ one small component
→ dependencies replaced

Integration Test
→ multiple real components work together
```

---

# Current Architecture

The project evolved from:

```text
Handler
↓
in-memory []Task
```

to:

```text
Handler
↓
Repository
↓
MySQL
```

and now:

```text
Handler
↓
Service
↓
Repository
↓
MySQL
```

Testing follows the same layers:

```text
Handler Test
↓
Fake Service
```

```text
Service Test
↓
Fake Repository
```

---

# Key Concepts

```text
Handler
→ HTTP concerns

Service
→ business logic

Repository
→ persistence/database logic

TaskRepository interface
→ separates Service from MySQL implementation

TaskService interface
→ separates Handler from Service implementation

?
→ MySQL parameter placeholder

ExecContext
→ INSERT / UPDATE / DELETE

QueryRowContext
→ one row

QueryContext
→ multiple rows

sql.ErrNoRows
→ database record not found

httptest
→ test HTTP handlers without running a server

table-driven test
→ same test logic with multiple cases

fake dependency
→ isolate one layer during testing

go test ./...
→ run all Go tests

go test -cover
→ measure statement coverage
```

---

# Key Takeaway

The main architecture is:

```text
HTTP
↓
Handler
↓
Service
↓
Repository
↓
database/sql
↓
MySQL
```

The testing architecture is:

```text
Handler
↓
Fake Service
```

and:

```text
Service
↓
Fake Repository
```

The main rule is:

> Keep HTTP, business logic, and database access in separate layers, and test each layer with controlled fake dependencies.