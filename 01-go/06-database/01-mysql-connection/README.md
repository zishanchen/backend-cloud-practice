# MySQL Connection with `database/sql`

A minimal Go example for connecting to MySQL using the standard `database/sql` package.

## Concepts Practiced

- `database/sql`
- MySQL database driver
- `sql.Open`
- `PingContext`
- `DB.Close`
- environment-variable DSN
- connection pools
- connection timeouts
- `DB.Stats`

## Connection Flow

```text
MYSQL_DSN
↓
sql.Open()
↓
*sql.DB
↓
configure connection pool
↓
PingContext()
↓
MySQL connection verified
↓
use database
↓
Close()
```

## Install the Driver

```bash
go get github.com/go-sql-driver/mysql@latest
go mod tidy
```

The driver is registered using a blank import:

```go
import _ "github.com/go-sql-driver/mysql"
```

This allows:

```go
sql.Open("mysql", dsn)
```

## DSN

The application reads the DSN from:

```text
MYSQL_DSN
```

Example:

```bash
export MYSQL_DSN='app:app@tcp(127.0.0.1:3306)/taskdb?parseTime=true'
```

Do not store real database credentials in source code.

## Run

From the repository root:

```bash
go run ./01-go/06-database/01-mysql-connection
```

Example output:

```text
MySQL connection OK
Connection pool:
  max open: 10
  open: 1
  in use: 0
  idle: 1
```

## `sql.Open`

`sql.Open` creates an `*sql.DB` database handle.

It does not necessarily establish a database connection immediately.

```go
db, err := sql.Open("mysql", dsn)
```

Think of `sql.DB` as:

```text
database handle
+
connection pool manager
```

It should normally be created once and reused.

## `PingContext`

Use `PingContext` to verify that the database is actually reachable.

```go
if err := db.PingContext(ctx); err != nil {
	...
}
```

A timeout context prevents the connection attempt from waiting indefinitely.

```go
ctx, cancel := context.WithTimeout(
	context.Background(),
	3*time.Second,
)
defer cancel()
```

## Connection Pool

`sql.DB` manages a pool of database connections automatically.

```text
Application
     ↓
   sql.DB
     ↓
Connection Pool
 ├── connection 1
 ├── connection 2
 └── connection 3
     ↓
   MySQL
```

Operations borrow connections from the pool and return them when finished.

### Maximum Open Connections

```go
db.SetMaxOpenConns(10)
```

Limits the total number of open database connections.

### Maximum Idle Connections

```go
db.SetMaxIdleConns(5)
```

Controls how many unused connections are kept ready for reuse.

### Maximum Connection Lifetime

```go
db.SetConnMaxLifetime(
	3 * time.Minute,
)
```

Limits how long one connection may be reused.

### Maximum Idle Time

```go
db.SetConnMaxIdleTime(
	time.Minute,
)
```

Limits how long an unused connection may remain idle.

## Pool Statistics

```go
stats := db.Stats()
```

Useful fields include:

```text
MaxOpenConnections
OpenConnections
InUse
Idle
```

Conceptually:

```text
OpenConnections
=
InUse + Idle
```

## Close

```go
db.Close()
```

closes the database handle and its resources.

For a long-running application, the normal lifecycle is:

```text
application starts
↓
open sql.DB once
↓
reuse connection pool
↓
application shuts down
↓
db.Close()
```

Do not open and close a new `sql.DB` for every HTTP request.

## Key Takeaways

```text
sql.Open
→ create database handle / pool

PingContext
→ verify real connectivity

sql.DB
→ reusable connection pool

environment variable
→ keep DSN out of source code

Close
→ release database resources at application shutdown
```