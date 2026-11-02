# Go Task HTTP Server

A small HTTP server built with Go's standard `net/http` package.

This exercise focuses on creating routes, handling HTTP methods, returning JSON responses, setting HTTP status codes, and testing endpoints with `curl`.

## Endpoints

The server provides two endpoints:

```text
GET /health
GET /tasks
```

### `GET /health`

Checks whether the service is running.

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

### `GET /tasks`

Returns a list of tasks as JSON.

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

The server starts on:

```text
http://localhost:8080
```

## Test with `curl`

### Check `/health`

```bash
curl -i http://localhost:8080/health
```

Example:

```text
HTTP/1.1 200 OK
Content-Type: application/json

{"status":"ok"}
```

### Check `/tasks`

```bash
curl -i http://localhost:8080/tasks
```

### Check only the HTTP status code

```bash
curl -s -o /dev/null -w "%{http_code}\n" \
  http://localhost:8080/health
```

Expected:

```text
200
```

### Test an unsupported HTTP method

```bash
curl -i -X POST http://localhost:8080/tasks
```

Expected:

```text
HTTP/1.1 405 Method Not Allowed
```

## Main Concepts Practiced

### `http.NewServeMux()`

Creates a request router:

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
HTTP Response
```

---

### HTTP Handlers

A handler has the form:

```go
func handler(w http.ResponseWriter, r *http.Request) {
}
```

`r` contains information about the incoming request, such as:

```text
HTTP method
URL
headers
body
```

`w` is used to create the response.

---

### HTTP Method Validation

The endpoints currently support only `GET`.

```go
if r.Method != http.MethodGet {
	w.WriteHeader(http.StatusMethodNotAllowed)
	return
}
```

This returns:

```text
405 Method Not Allowed
```

for unsupported methods.

---

### JSON Responses

The response content type is set with:

```go
w.Header().Set("Content-Type", "application/json")
```

Go values are encoded as JSON using:

```go
json.NewEncoder(w).Encode(...)
```

For example:

```go
json.NewEncoder(w).Encode(tasks)
```

---

### HTTP Status Codes

The server uses Go HTTP constants instead of raw numbers:

```go
http.StatusOK
http.StatusMethodNotAllowed
```

Equivalent values:

```text
200 → OK
405 → Method Not Allowed
```

The response order should normally be:

```text
Set headers
↓
Write status
↓
Write body
```

Example:

```go
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(http.StatusOK)
json.NewEncoder(w).Encode(tasks)
```

## Key Pattern

For a basic Go HTTP endpoint:

```text
Route
↓
Handler
↓
Check HTTP Method
↓
Set Headers
↓
Set Status Code
↓
Encode JSON Response
```

For this project:

```text
GET /health
→ 200
→ JSON health status

GET /tasks
→ 200
→ JSON task list

POST /tasks
→ 405
```

## What I Learned

This exercise connects several earlier Go topics:

```text
structs
+ slices
+ JSON
+ functions
+ error handling
+ net/http
+ HTTP methods
+ status codes
+ curl
```

The main request flow to remember is:

```text
curl
→ HTTP request
→ ServeMux
→ handler
→ ResponseWriter
→ status + headers + JSON
```