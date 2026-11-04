# Go Task HTTP Server

A small REST-style HTTP server built with Go's standard `net/http` and `encoding/json` packages.

The project started with simple HTTP handlers and was gradually extended with:

- `GET /health`
- `GET /tasks`
- `POST /tasks`
- `GET /tasks/{id}`
- JSON decoding and encoding
- Input validation
- Structured error responses
- Go 1.22 `ServeMux` method-based routing
- Go 1.22 path parameters with `Request.PathValue()`

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

Checks whether the service is running.

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

---

## `GET /tasks`

Returns all tasks.

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
    "title": "Learn Go 1.22 ServeMux",
    "done": false
  }
]
```

Expected status:

```text
200 OK
```

---

## `GET /tasks/{id}`

Returns one task by ID.

Example:

```bash
curl -i http://localhost:8080/tasks/2
```

Example response:

```json
{
  "id": 2,
  "title": "Build /health",
  "done": true
}
```

Expected status:

```text
200 OK
```

If the task does not exist:

```bash
curl -i http://localhost:8080/tasks/99
```

Response:

```json
{
  "error": "task not found"
}
```

Expected status:

```text
404 Not Found
```

If the path parameter is not a valid integer:

```bash
curl -i http://localhost:8080/tasks/abc
```

Response:

```json
{
  "error": "invalid task id"
}
```

Expected status:

```text
400 Bad Request
```

---

## `POST /tasks`

Creates a new task.

```bash
curl -i \
  -X POST \
  http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Practice ServeMux routing","done":false}'
```

Example response:

```json
{
  "id": 4,
  "title": "Practice ServeMux routing",
  "done": false
}
```

Expected status:

```text
201 Created
```

---

# Validation

The task title is required.

For example:

```bash
curl -i \
  -X POST \
  http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":""}'
```

Response:

```json
{
  "error": "title is required"
}
```

Expected status:

```text
400 Bad Request
```

Whitespace is removed with:

```go
strings.TrimSpace(input.Title)
```

so this is also invalid:

```json
{
  "title": "   "
}
```

---

# Invalid JSON

Malformed JSON returns:

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

Response:

```json
{
  "error": "invalid JSON body"
}
```

The decoder also rejects unknown fields:

```go
decoder.DisallowUnknownFields()
```

For example:

```json
{
  "title": "Learn Go",
  "done": false,
  "priority": 5
}
```

is rejected because `priority` is not defined in `CreateTaskRequest`.

---

# Go 1.22 ServeMux Routing

The most important change in this version is the routing style.

Previously, the route was registered only by path:

```go
mux.HandleFunc("/tasks", tasksHandler)
```

and the handler had to inspect the method:

```go
switch r.Method {
case http.MethodGet:
	getTasks(w)

case http.MethodPost:
	createTask(w, r)
}
```

With Go 1.22-style `ServeMux` routing, the HTTP method can be part of the route pattern:

```go
mux.HandleFunc(
	"GET /health",
	healthHandler,
)

mux.HandleFunc(
	"GET /tasks",
	getTasksHandler,
)

mux.HandleFunc(
	"POST /tasks",
	createTaskHandler,
)
```

The router now decides based on both:

```text
HTTP Method
+
Path
```

Example:

```text
GET /tasks
→ getTasksHandler

POST /tasks
→ createTaskHandler
```

This removes the need for a manual `switch r.Method` inside the handler.

---

# Path Parameters

Go 1.22 `ServeMux` also supports path wildcards.

The route:

```go
mux.HandleFunc(
	"GET /tasks/{id}",
	getTaskHandler,
)
```

matches requests such as:

```text
GET /tasks/1
GET /tasks/2
GET /tasks/10
```

The path value is read with:

```go
idText := r.PathValue("id")
```

For:

```text
GET /tasks/2
```

the value is:

```text
"2"
```

It can then be converted into an integer:

```go
id, err := strconv.Atoi(idText)
```

Request flow:

```text
GET /tasks/2
↓
ServeMux matches /tasks/{id}
↓
getTaskHandler
↓
r.PathValue("id")
↓
"2"
↓
strconv.Atoi()
↓
2
```

---

# Current Routes

```text
GET  /health
GET  /tasks
GET  /tasks/{id}
POST /tasks
```

Routing flow:

```text
Request
↓
ServeMux
↓
Match Method + Path
↓
Handler
```

Example:

```text
GET /tasks
→ getTasksHandler

GET /tasks/2
→ getTaskHandler

POST /tasks
→ createTaskHandler
```

---

# JSON Decoder

For `POST /tasks`, JSON from the request body is converted into a Go struct.

```go
decoder := json.NewDecoder(r.Body)

if err := decoder.Decode(&input); err != nil {
	// return 400
}
```

Flow:

```text
HTTP request body
↓
r.Body
↓
json.Decoder
↓
Go struct
```

Example input:

```json
{
  "title": "Practice Go",
  "done": false
}
```

becomes:

```go
CreateTaskRequest{
	Title: "Practice Go",
	Done:  false,
}
```

A pointer is passed because `Decode()` must modify the struct:

```go
decoder.Decode(&input)
```

---

# JSON Encoder

Responses are converted from Go values to JSON with:

```go
json.NewEncoder(w).Encode(data)
```

Flow:

```text
Go value
↓
json.Encoder
↓
JSON
↓
http.ResponseWriter
↓
Client
```

The rule to remember is:

```text
Decoder = JSON → Go
Encoder = Go → JSON
```

---

# Response Helpers

## `writeJSON`

The helper:

```go
func writeJSON(
	w http.ResponseWriter,
	status int,
	data any,
)
```

handles:

```text
Content-Type
+
status code
+
JSON body
```

Instead of repeating:

```go
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(...)
json.NewEncoder(w).Encode(...)
```

handlers can use:

```go
writeJSON(
	w,
	http.StatusOK,
	tasks,
)
```

---

## `writeError`

Errors use a common response format:

```go
type ErrorResponse struct {
	Error string `json:"error"`
}
```

Example:

```json
{
  "error": "task not found"
}
```

The helper:

```go
writeError(
	w,
	http.StatusNotFound,
	"task not found",
)
```

keeps error responses consistent.

---

# HTTP Status Codes Used

```text
200 OK
201 Created
400 Bad Request
404 Not Found
405 Method Not Allowed
```

Corresponding Go constants include:

```go
http.StatusOK
http.StatusCreated
http.StatusBadRequest
http.StatusNotFound
http.StatusMethodNotAllowed
```

---

# Response Order

The normal response order is:

```text
1. Set headers
2. Write status code
3. Write body
```

Example:

```go
w.Header().Set(
	"Content-Type",
	"application/json",
)

w.WriteHeader(http.StatusOK)

json.NewEncoder(w).Encode(data)
```

The `writeJSON()` helper performs these steps in the correct order.

---

# Request Flow Examples

## `GET /tasks/{id}`

```text
curl
↓
GET /tasks/2
↓
ServeMux
↓
"GET /tasks/{id}"
↓
getTaskHandler
↓
PathValue("id")
↓
Find task
↓
Encoder
↓
200 + JSON
```

## `POST /tasks`

```text
curl
↓
POST /tasks
↓
ServeMux
↓
createTaskHandler
↓
Request Body
↓
Decoder
↓
CreateTaskRequest
↓
Validation
↓
Create Task
↓
Encoder
↓
201 + JSON
```

---

# Test with `curl`

Start the server:

```bash
go run ./04-http/task-server
```

Health check:

```bash
curl -i http://localhost:8080/health
```

Get all tasks:

```bash
curl -i http://localhost:8080/tasks
```

Get one task:

```bash
curl -i http://localhost:8080/tasks/2
```

Get missing task:

```bash
curl -i http://localhost:8080/tasks/99
```

Invalid ID:

```bash
curl -i http://localhost:8080/tasks/abc
```

Create task:

```bash
curl -i \
  -X POST \
  http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"Practice ServeMux","done":false}'
```

Invalid title:

```bash
curl -i \
  -X POST \
  http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":""}'
```

Invalid JSON:

```bash
curl -i \
  -X POST \
  http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":}'
```

Unsupported method:

```bash
curl -i \
  -X DELETE \
  http://localhost:8080/tasks
```

---

# Concepts Practiced

This version combines:

```text
net/http
encoding/json
ServeMux
method-based routing
path parameters
Request.PathValue()
strconv.Atoi()
structs
slices
JSON tags
validation
HTTP status codes
error responses
curl
```

The project now demonstrates a more structured HTTP flow:

```text
Route
↓
Handler
↓
Parse Input
↓
Validate
↓
Application Logic
↓
Encode Response
```

---

# Evolution of the Project

The HTTP server evolved in several steps:

```text
Step 1
Basic net/http server
↓
GET /health
GET /tasks

Step 2
JSON Decoder / Encoder
↓
POST /tasks
↓
Validation
↓
Structured errors

Step 3
Go 1.22 ServeMux
↓
Method-based routes
↓
GET /tasks/{id}
↓
PathValue("id")
```

The main routing improvement is:

```text
Before

"/tasks"
↓
tasksHandler
↓
switch r.Method
```

becoming:

```text
Now

"GET /tasks"
→ getTasksHandler

"POST /tasks"
→ createTaskHandler

"GET /tasks/{id}"
→