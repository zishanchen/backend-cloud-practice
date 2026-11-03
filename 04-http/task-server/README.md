# Go Task HTTP Server

A small REST-style HTTP server built with Go's standard `net/http` and `encoding/json` packages.

The project started with simple `GET` endpoints and was extended with `POST /tasks` to practice JSON decoding, validation, HTTP status codes, error responses, and response encoding.

## Features

- `GET /health`
- `GET /tasks`
- `POST /tasks`
- JSON request decoding
- JSON response encoding
- Input validation
- Structured JSON error responses
- HTTP status handling
- Testing endpoints with `curl`

## Project Structure

```text
task-server/
├── main.go
└── README.md
```

## Run the Server

From the repository root:

```bash
go run ./04-http/task-server
```

The server starts at:

```text
http://localhost:8080
```

---

# API Endpoints

## `GET /health`

Checks whether the HTTP service is running.

Request:

```bash
curl -i http://localhost:8080/health
```

Example response:

```text
HTTP/1.1 200 OK
Content-Type: application/json
```

```json
{
  "status": "ok"
}
```

Expected status:

```text
200 OK
```

---

## `GET /tasks`

Returns all tasks as JSON.

Request:

```bash
curl -i http://localhost:8080/tasks
```

Example response:

```json
[
  {
    "id": 1,
    "title": "Learn net/http",
    "done": true
  },
  {
    "id": 2,
    "title": "Build /health",
    "done": true
  },
  {
    "id": 3,
    "title": "Build /tasks",
    "done": false
  }
]
```

Expected status:

```text
200 OK
```

---

## `POST /tasks`

Creates a new task.

Request:

```bash
curl -i \
  -X POST \
  http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Practice SQL SELECT","done":false}'
```

Example response:

```text
HTTP/1.1 201 Created
Content-Type: application/json
```

```json
{
  "id": 4,
  "title": "Practice SQL SELECT",
  "done": false
}
```

Expected status:

```text
201 Created
```

After creating the task, calling:

```bash
curl http://localhost:8080/tasks
```

returns the updated in-memory task list.

---

# Validation

The server validates the incoming task title.

For example:

```bash
curl -i \
  -X POST \
  http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":""}'
```

Response:

```text
HTTP/1.1 400 Bad Request
```

```json
{
  "error": "title is required"
}
```

Whitespace-only titles are also rejected.

For example:

```json
{
  "title": "   "
}
```

is treated as an empty title after `strings.TrimSpace()`.

---

# Invalid JSON

Malformed JSON also returns:

```text
400 Bad Request
```

Example:

```bash
curl -i \
  -X POST \
  http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":}'
```

Example response:

```json
{
  "error": "invalid JSON body"
}
```

---

# Unsupported HTTP Methods

Only the configured HTTP methods are accepted.

For example:

```bash
curl -i \
  -X DELETE \
  http://localhost:8080/tasks
```

returns:

```text
405 Method Not Allowed
```

with a JSON error response:

```json
{
  "error": "method not allowed"
}
```

---

# Main Data Structures

## `Task`

```go
type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}
```

A `Task` represents the full task returned by the API.

Example:

```json
{
  "id": 4,
  "title": "Practice SQL SELECT",
  "done": false
}
```

## `CreateTaskRequest`

```go
type CreateTaskRequest struct {
	Title string `json:"title"`
	Done  bool   `json:"done"`
}
```

This struct represents the JSON input required to create a task.

The client does not provide the task ID.

The server generates it.

## `ErrorResponse`

```go
type ErrorResponse struct {
	Error string `json:"error"`
}
```

Used to return consistent JSON error messages.

Example:

```json
{
  "error": "title is required"
}
```

---

# JSON Decoder

For `POST /tasks`, the request body contains JSON.

The server converts it into a Go struct using:

```go
decoder := json.NewDecoder(r.Body)

if err := decoder.Decode(&input); err != nil {
	// return 400
}
```

The flow is:

```text
HTTP request body
        ↓
      r.Body
        ↓
 json.Decoder
        ↓
     Go struct
```

For example:

```json
{
  "title": "Practice SQL",
  "done": false
}
```

becomes approximately:

```go
CreateTaskRequest{
	Title: "Practice SQL",
	Done:  false,
}
```

Because `Decode()` modifies the target value, a pointer is passed:

```go
decoder.Decode(&input)
```

---

# Unknown JSON Fields

The decoder is configured with:

```go
decoder.DisallowUnknownFields()
```

This makes unexpected JSON fields invalid.

For example:

```json
{
  "title": "Learn Go",
  "done": false,
  "priority": 5
}
```

will be rejected because `priority` is not part of `CreateTaskRequest`.

---

# JSON Encoder

Responses are written using:

```go
json.NewEncoder(w).Encode(data)
```

The flow is:

```text
Go value
   ↓
json.Encoder
   ↓
JSON
   ↓
http.ResponseWriter
   ↓
HTTP client
```

A useful rule to remember is:

```text
Decoder = JSON → Go

Encoder = Go → JSON
```

---

# `writeJSON`

To avoid repeating:

```go
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(...)
json.NewEncoder(w).Encode(...)
```

the project uses a helper function:

```go
func writeJSON(
	w http.ResponseWriter,
	status int,
	data any,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(data)
}
```

This makes handlers shorter and keeps JSON responses consistent.

Example:

```go
writeJSON(
	w,
	http.StatusOK,
	tasks,
)
```

---

# `writeError`

Error responses use another helper:

```go
func writeError(
	w http.ResponseWriter,
	status int,
	message string,
) {
	writeJSON(
		w,
		status,
		ErrorResponse{
			Error: message,
		},
	)
}
```

Instead of repeating error response logic, handlers can call:

```go
writeError(
	w,
	http.StatusBadRequest,
	"title is required",
)
```

---

# Routing

The server uses:

```go
mux := http.NewServeMux()
```

Routes are registered with:

```go
mux.HandleFunc("/health", healthHandler)
mux.HandleFunc("/tasks", tasksHandler)
```

Request flow:

```text
Client
↓
HTTP Request
↓
ServeMux
↓
Handler
↓
Application Logic
↓
HTTP Response
```

---

# Handling Multiple Methods

The `/tasks` endpoint supports both:

```text
GET
POST
```

The handler uses:

```go
switch r.Method {
case http.MethodGet:
	getTasks(w)

case http.MethodPost:
	createTask(w, r)

default:
	writeError(
		w,
		http.StatusMethodNotAllowed,
		"method not allowed",
	)
}
```

This separates:

```text
GET /tasks
→ read tasks

POST /tasks
→ create task
```

---

# HTTP Status Codes Used

```text
200 OK
201 Created
400 Bad Request
405 Method Not Allowed
```

Go constants:

```go
http.StatusOK
http.StatusCreated
http.StatusBadRequest
http.StatusMethodNotAllowed
```

Using named constants is clearer than writing raw numbers such as:

```go
w.WriteHeader(201)
```

---

# Important Response Order

The response should normally be constructed in this order:

```text
1. Set headers
2. Write status code
3. Write response body
```

Example:

```go
w.Header().Set(
	"Content-Type",
	"application/json",
)

w.WriteHeader(http.StatusCreated)

json.NewEncoder(w).Encode(task)
```

Once the response body or status has been written, changing headers may be too late.

---

# Request Flow for `POST /tasks`

The complete flow is:

```text
Client sends POST JSON
        ↓
POST /tasks
        ↓
tasksHandler
        ↓
createTask
        ↓
json.Decoder
        ↓
CreateTaskRequest
        ↓
Validate title
        ↓
Generate ID
        ↓
Append to tasks slice
        ↓
json.Encoder
        ↓
201 Created + JSON
```

This is the main pattern practiced in this exercise:

```text
POST JSON
→ Decode
→ Validate
→ Process
→ Encode
→ HTTP Response
```

---

# Test with `curl`

## Health Check

```bash
curl -i http://localhost:8080/health
```

## Get Tasks

```bash
curl -i http://localhost:8080/tasks
```

## Create Task

```bash
curl -i \
  -X POST \
  http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Practice SQL","done":false}'
```

## Empty Title

```bash
curl -i \
  -X POST \
  http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":""}'
```

## Invalid JSON

```bash
curl -i \
  -X POST \
  http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":}'
```

## Unsupported Method

```bash
curl -i \
  -X DELETE \
  http://localhost:8080/tasks
```

## Only Show Status Code

```bash
curl -s \
  -o /dev/null \
  -w "%{http_code}\n" \
  http://localhost:8080/tasks
```

---

# Concepts Practiced

This project combines several Go topics:

```text
net/http
encoding/json
structs
slices
JSON tags
functions
pointers
error handling
HTTP methods
HTTP status codes
request bodies
response bodies
validation
curl
```

It also introduces a basic separation of concerns:

```text
Routing
↓
Handler
↓
Validation
↓
Application Logic
↓
Response
```

---

# Current Limitations

The project intentionally remains simple.

Current limitations include:

- Tasks are stored only in memory.
- Data disappears when the server restarts.
- There is no database.
- There is no concurrency protection around the global task slice.
- There is no authentication.
- There are no update or delete endpoints yet.
- There are no automated HTTP handler tests yet.

These can be added in later exercises.

---

# Key Takeaways

The most important concepts from this version are:

```text
GET
→ return data
→ Encoder
→ JSON

POST
→ receive JSON
→ Decoder
→ validate
→ create data
→ Encoder
→ JSON
```

And the main rule to remember:

```text
Decoder = JSON → Go
Encoder = Go → JSON
```

The complete request lifecycle is:

```text
curl
→ HTTP request
→ ServeMux
→ Handler
→ Decode
→ Validate
→ Process
→ Encode
→ Status + JSON response
```