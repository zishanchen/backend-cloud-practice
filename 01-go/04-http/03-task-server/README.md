# Go Task HTTP Server

A small REST-style HTTP server built with Go's standard library.

The project gradually evolved from basic `net/http` handlers into a more structured in-memory HTTP API with:

- Go 1.22 `ServeMux` routing
- 5 HTTP endpoints
- JSON request and response handling
- Input validation
- Path parameters
- Structured error responses
- Request ID middleware
- Logging middleware
- HTTP status tracking
- Basic request tracing

## Project Structure

```text
03-task-server/
├── main.go
├── middleware.go
└── README.md
```

## Run the Server

From the repository root:

```bash
go run ./01-go/04-http/03-task-server
```

The server starts at:

```text
http://localhost:8080
```

---

# API Endpoints

The current in-memory API provides 5 endpoints:

```text
GET     /health
GET     /tasks
GET     /tasks/{id}
POST    /tasks
DELETE  /tasks/{id}
```

---

## `GET /health`

Checks whether the service is running.

```bash
curl -i http://localhost:8080/health
```

Example response:

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

If the ID is invalid:

```bash
curl -i http://localhost:8080/tasks/abc
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
  -d '{"title":"Practice middleware","done":false}'
```

Example response:

```json
{
  "id": 4,
  "title": "Practice middleware",
  "done": false
}
```

Expected status:

```text
201 Created
```

Invalid request:

```bash
curl -i \
  -X POST \
  http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":""}'
```

Expected status:

```text
400 Bad Request
```

---

## `DELETE /tasks/{id}`

Deletes a task by ID.

```bash
curl -i \
  -X DELETE \
  http://localhost:8080/tasks/2
```

Expected status:

```text
204 No Content
```

A successful delete returns no response body.

Verify that the task was removed:

```bash
curl -i http://localhost:8080/tasks/2
```

Expected result:

```text
404 Not Found
```

If the task does not exist:

```bash
curl -i \
  -X DELETE \
  http://localhost:8080/tasks/99
```

Expected status:

```text
404 Not Found
```

If the ID is invalid:

```bash
curl -i \
  -X DELETE \
  http://localhost:8080/tasks/abc
```

Expected status:

```text
400 Bad Request
```

---

# Go 1.22 ServeMux Routing

Routes are registered using both the HTTP method and path:

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
	"GET /tasks/{id}",
	getTaskHandler,
)

mux.HandleFunc(
	"POST /tasks",
	createTaskHandler,
)

mux.HandleFunc(
	"DELETE /tasks/{id}",
	deleteTaskHandler,
)
```

Routing is based on:

```text
HTTP Method
+
Path
```

For example:

```text
GET /tasks
→ getTasksHandler

POST /tasks
→ createTaskHandler

GET /tasks/2
→ getTaskHandler

DELETE /tasks/2
→ deleteTaskHandler
```

This avoids manually switching on:

```go
switch r.Method {
...
}
```

inside a shared handler.

---

# Path Parameters

Both:

```text
GET /tasks/{id}
DELETE /tasks/{id}
```

use a path parameter.

The value is read with:

```go
idText := r.PathValue("id")
```

For:

```text
DELETE /tasks/2
```

the value is:

```text
"2"
```

It is converted to an integer with:

```go
id, err := strconv.Atoi(idText)
```

Flow:

```text
DELETE /tasks/2
↓
ServeMux
↓
DELETE /tasks/{id}
↓
PathValue("id")
↓
"2"
↓
strconv.Atoi()
↓
2
```

---

# In-Memory Storage

Tasks are currently stored in a Go slice:

```go
var tasks = []Task{
	...
}
```

This means the API does not yet use a database.

Operations modify the slice directly:

```text
GET
→ read tasks

POST
→ append task

DELETE
→ remove task from slice
```

Example delete operation:

```go
tasks = append(
	tasks[:i],
	tasks[i+1:]...,
)
```

This joins:

```text
everything before index i
+
everything after index i
```

and removes the matching task.

The data exists only while the Go process is running.

Restarting the server resets the in-memory data.

---

# JSON Decoder

For `POST /tasks`, the request body contains JSON.

```go
decoder := json.NewDecoder(r.Body)
decoder.DisallowUnknownFields()

if err := decoder.Decode(&input); err != nil {
	...
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

Example:

```json
{
  "title": "Practice Go",
  "done": false
}
```

becomes approximately:

```go
CreateTaskRequest{
	Title: "Practice Go",
	Done:  false,
}
```

Remember:

```text
Decoder = JSON → Go
```

---

# JSON Encoder

Responses are written with:

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
HTTP response
```

Remember:

```text
Encoder = Go → JSON
```

---

# Validation

The task title must not be empty.

```go
input.Title = strings.TrimSpace(input.Title)

if input.Title == "" {
	writeError(
		w,
		http.StatusBadRequest,
		"title is required",
	)
	return
}
```

Task IDs are also validated:

```go
id, err := strconv.Atoi(idText)

if err != nil {
	writeError(
		w,
		http.StatusBadRequest,
		"invalid task id",
	)
	return
}
```

---

# Response Helpers

## `writeJSON`

The helper centralizes:

```text
Content-Type
+
HTTP status
+
JSON encoding
```

Example:

```go
writeJSON(
	w,
	http.StatusOK,
	tasks,
)
```

## `writeError`

Errors use a consistent JSON format:

```json
{
  "error": "task not found"
}
```

Example:

```go
writeError(
	w,
	http.StatusNotFound,
	"task not found",
)
```

---

# Middleware

Middleware is reusable logic that wraps an `http.Handler`.

The core pattern is:

```go
func middleware(
	next http.Handler,
) http.Handler {

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// before

			next.ServeHTTP(w, r)

			// after
		},
	)
}
```

The most important line is:

```go
next.ServeHTTP(w, r)
```

It passes the request to the next handler in the chain.

Typical middleware use cases include:

```text
logging
request IDs
authentication
authorization
metrics
CORS
rate limiting
panic recovery
```

---

# Middleware Chain

The task server wraps the router like this:

```go
handler := requestIDMiddleware(
	loggingMiddleware(
		mux,
	),
)
```

Structure:

```text
requestIDMiddleware
└── loggingMiddleware
    └── ServeMux
```

Request flow:

```text
Client
↓
requestIDMiddleware
↓
loggingMiddleware
↓
ServeMux
↓
Route Handler
↓
Response
```

---

# Request ID Middleware

Each request receives a request ID.

The middleware first checks:

```go
r.Header.Get("X-Request-ID")
```

If the client did not provide one, the server generates one.

Example:

```text
req-1
req-2
req-3
```

The request ID is also returned in the response header:

```text
X-Request-ID: req-1
```

This allows correlation between:

```text
request
↔
server logs
↔
response
```

---

# Request Context

The request ID is stored in the request context:

```go
ctx := context.WithValue(
	r.Context(),
	requestIDKey,
	requestID,
)
```

The modified request is passed downstream:

```go
next.ServeHTTP(
	w,
	r.WithContext(ctx),
)
```

Later middleware or handlers can retrieve the same request ID.

Typical request-specific context values include:

```text
request ID
trace ID
authenticated user
deadline
cancellation
```

---

# Logging Middleware

The logging middleware measures request duration and records request information.

Example:

```text
request_id=req-3 method=GET path=/tasks/2 status=200 duration=95µs
```

A DELETE request might produce:

```text
request_id=req-4 method=DELETE path=/tasks/2 status=204 duration=80µs
```

Fields:

```text
request_id
→ identifies the request

method
→ GET / POST / DELETE / ...

path
→ requested endpoint

status
→ HTTP response status

duration
→ request processing time
```

Basic flow:

```text
record start time
↓
next.ServeHTTP()
↓
handler executes
↓
record status
↓
calculate duration
↓
write log
```

---

# Response Recorder

The logging middleware needs to know which status code the handler returned.

`http.ResponseWriter` is wrapped:

```go
type responseRecorder struct {
	http.ResponseWriter
	status int
}
```

The wrapper intercepts:

```go
func (rec *responseRecorder) WriteHeader(status int) {
	rec.status = status
	rec.ResponseWriter.WriteHeader(status)
}
```

This allows the logging middleware to record statuses such as:

```text
200
201
204
400
404
405
```

---

# Why Middleware Order Matters

Current order:

```text
request ID
↓
logging
↓
router
```

The request ID middleware runs first, so the logging middleware can read the generated request ID.

```text
request
↓
requestIDMiddleware
    create request ID
↓
loggingMiddleware
    read request ID
↓
handler
```

Middleware order can change behavior.

---

# Test Middleware

Start the server:

```bash
go run ./01-go/04-http/03-task-server
```

Normal request:

```bash
curl -i http://localhost:8080/tasks
```

Response header:

```text
X-Request-Id: req-1
```

Test DELETE:

```bash
curl -i \
  -X DELETE \
  http://localhost:8080/tasks/2
```

Server log:

```text
request_id=req-2 method=DELETE path=/tasks/2 status=204 duration=...
```

---

# Custom Request ID

A client can provide its own request ID:

```bash
curl -i \
  -H "X-Request-ID: test-123" \
  http://localhost:8080/tasks
```

Response:

```text
X-Request-Id: test-123
```

Log:

```text
request_id=test-123 method=GET path=/tasks status=200 duration=...
```

---

# HTTP Status Codes Used

```text
200 OK
201 Created
204 No Content
400 Bad Request
404 Not Found
405 Method Not Allowed
```

Go constants include:

```go
http.StatusOK
http.StatusCreated
http.StatusNoContent
http.StatusBadRequest
http.StatusNotFound
http.StatusMethodNotAllowed
```

---

# Full Request Lifecycles

## `GET /tasks/2`

```text
curl
↓
HTTP request
↓
requestIDMiddleware
↓
loggingMiddleware
↓
ServeMux
↓
GET /tasks/{id}
↓
getTaskHandler
↓
PathValue("id")
↓
application logic
↓
writeJSON
↓
responseRecorder
↓
logging
↓
HTTP response
```

## `POST /tasks`

```text
curl
↓
POST /tasks
↓
request ID
↓
logging
↓
ServeMux
↓
createTaskHandler
↓
JSON Decoder
↓
validation
↓
create task
↓
JSON Encoder
↓
201 Created
↓
logging
↓
response
```

## `DELETE /tasks/2`

```text
curl
↓
DELETE /tasks/2
↓
requestIDMiddleware
↓
loggingMiddleware
↓
ServeMux
↓
DELETE /tasks/{id}
↓
deleteTaskHandler
↓
PathValue("id")
↓
find matching task
↓
remove task from slice
↓
204 No Content
↓
logging
↓
response
```

---

# Concepts Practiced

This project currently practices:

```text
net/http
ServeMux
Go 1.22 routing
HTTP methods
GET
POST
DELETE
path parameters
Request.PathValue()
encoding/json
JSON Decoder
JSON Encoder
validation
HTTP status codes
structured errors
in-memory CRUD-style operations
slice deletion
http.Handler
http.HandlerFunc
middleware
handler wrapping
request context
request IDs
logging
response writer wrapping
request duration
curl
```

---

# Project Evolution

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
validation
↓
structured errors

Step 3
Go 1.22 ServeMux
↓
method-based routes
↓
GET /tasks/{id}
↓
PathValue()

Step 4
Middleware
↓
request ID
↓
logging
↓
status tracking
↓
request duration

Step 5
In-memory API expansion
↓
DELETE /tasks/{id}
↓
5 endpoints
↓
204 No Content
↓
slice deletion
```

---

# Current Limitations

The project intentionally remains simple.

Current limitations:

- Tasks are stored only in memory.
- Data disappears when the server restarts.
- PostgreSQL is not connected to the Go API yet.
- The global task slice has no concurrency protection.
- There is no authentication.
- There is no update endpoint such as `PUT` or `PATCH`.
- There are no automated HTTP handler tests yet.
- Request IDs are locally generated rather than using a distributed tracing system.

These can be introduced in later exercises.

---

# Key Takeaways

## Routing

```text
METHOD + PATH
→ Handler
```

Examples:

```text
GET /tasks
→ getTasksHandler

POST /tasks
→ createTaskHandler

GET /tasks/{id}
→ getTaskHandler

DELETE /tasks/{id}
→ deleteTaskHandler
```

## JSON

```text
Decoder = JSON → Go

Encoder = Go → JSON
```

## In-Memory API

```text
GET
→ read

POST
→ create

DELETE
→ remove
```

## Middleware

```text
Middleware receives a Handler
↓
wraps it
↓
returns another Handler
```

Core signature:

```go
func(http.Handler) http.Handler
```

Core call:

```go
next.ServeHTTP(w, r)
```

Current request flow:

```text
request
↓
request ID middleware
↓
logging middleware
↓
ServeMux
↓
handler
↓
response
```

The main idea to remember is:

> Middleware handles reusable request-level concerns, while route handlers contain endpoint-specific application logic.