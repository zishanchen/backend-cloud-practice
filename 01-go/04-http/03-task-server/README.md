# Go Task HTTP Server

A small REST-style HTTP server built with Go's standard library.

The project gradually evolved from basic `net/http` handlers into a more structured HTTP service with:

- Go 1.22 `ServeMux` routing
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
```

This means routing is based on:

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
```

This removes the need to manually write:

```go
switch r.Method {
...
}
```

inside a shared handler.

---

# Path Parameters

The route:

```go
"GET /tasks/{id}"
```

contains a path wildcard.

The value is read with:

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

It is converted to an integer with:

```go
id, err := strconv.Atoi(idText)
```

Request flow:

```text
GET /tasks/2
↓
ServeMux
↓
GET /tasks/{id}
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

Example invalid request:

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

---

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

A middleware can run logic:

```text
before the handler
+
after the handler
```

Typical use cases include:

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

The task server currently wraps the router like this:

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

This helps correlate:

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

This is useful because request-specific data should not normally be stored in global variables.

Examples of request-specific context data:

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

Example log:

```text
request_id=req-3 method=GET path=/tasks/2 status=200 duration=95µs
```

Fields:

```text
request_id
→ identifies the request

method
→ GET / POST / ...

path
→ requested endpoint

status
→ HTTP response status

duration
→ request processing time
```

The basic flow is:

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

For that reason, `http.ResponseWriter` is wrapped:

```go
type responseRecorder struct {
	http.ResponseWriter
	status int
}
```

The wrapper intercepts:

```go
WriteHeader(status)
```

Example:

```go
func (rec *responseRecorder) WriteHeader(status int) {
	rec.status = status
	rec.ResponseWriter.WriteHeader(status)
}
```

This allows the logging middleware to record:

```text
200
201
400
404
```

after the handler finishes.

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

This allows the logging middleware to access the request ID that was created earlier.

Conceptually:

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

Middleware order can change behavior, so wrapper order matters.

---

# Test Middleware

Start the server:

```bash
go run ./01-go/04-http/03-task-server
```

Test normal request:

```bash
curl -i http://localhost:8080/tasks
```

The response should contain:

```text
X-Request-Id: req-1
```

Test another request:

```bash
curl -i http://localhost:8080/tasks/2
```

The server log should contain something similar to:

```text
request_id=req-2 method=GET path=/tasks/2 status=200 duration=...
```

---

# Custom Request ID

A client can provide its own request ID:

```bash
curl -i \
  -H "X-Request-ID: test-123" \
  http://localhost:8080/tasks
```

The response should contain:

```text
X-Request-Id: test-123
```

The server log should contain:

```text
request_id=test-123 method=GET path=/tasks status=200 duration=...
```

---

# HTTP Status Codes Used

```text
200 OK
201 Created
400 Bad Request
404 Not Found
405 Method Not Allowed
```

Go constants include:

```go
http.StatusOK
http.StatusCreated
http.StatusBadRequest
http.StatusNotFound
http.StatusMethodNotAllowed
```

---

# Full Request Lifecycle

A `GET /tasks/2` request now follows:

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

For `POST /tasks`:

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

---

# Concepts Practiced

This project currently practices:

```text
net/http
ServeMux
Go 1.22 routing
HTTP methods
path parameters
Request.PathValue()
encoding/json
JSON Decoder
JSON Encoder
validation
HTTP status codes
structured errors
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
```

---

# Current Limitations

The project intentionally remains simple.

Current limitations:

- Tasks are stored only in memory.
- Data disappears when the server restarts.
- There is no database integration yet.
- The global task slice has no concurrency protection.
- There is no authentication.
- There are no update or delete endpoints.
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

Example:

```text
GET /tasks
→ getTasksHandler

POST /tasks
→ createTaskHandler

GET /tasks/{id}
→ getTaskHandler
```

## JSON

```text
Decoder = JSON → Go

Encoder = Go → JSON
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

> Middleware is reusable logic that wraps HTTP handlers and runs before and/or after the actual handler.