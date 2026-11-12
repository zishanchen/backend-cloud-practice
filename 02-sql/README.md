# SQL Practice

This folder contains SQL exercises, schema design practice, sample data, JOIN exercises, and database schemas used with SQLite, PostgreSQL, and MySQL.

## Files

- `setup.sql` — creates sample data for basic SQLite SQL practice
- `select-basics.sql` — practices basic `SELECT`, `WHERE`, `ORDER BY`, and `LIMIT`
- `task-system-schema.sql` — defines the PostgreSQL task management database schema
- `task-system-seed.sql` — inserts sample data into the task management tables
- `task-system-queries.sql` — contains 10 SQL queries for filtering, joins, grouping, and aggregation
- `users-tasks-join.sql` — practices joining `users` and `tasks`
- `mysql-task-api-schema.sql` — defines the MySQL table used by the Go repository practice

---

# Task System Schema

The PostgreSQL task system contains four tables:

```text
users
tasks
status_history
audit_logs
```

## `users`

Stores application users.

Important fields:

```text
id
name
email
```

## `tasks`

Stores the current state of each task.

Important fields:

```text
id
title
status
user_id
created_at
```

Relationship:

```text
users.id
   ↑
   │ FOREIGN KEY
   │
tasks.user_id
```

This means:

> Each task belongs to one user.

## `status_history`

Stores task status changes.

Important fields:

```text
id
task_id
old_status
new_status
changed_at
```

Relationship:

```text
tasks.id
   ↑
   │ FOREIGN KEY
   │
status_history.task_id
```

This means:

> Each status history record belongs to one task.

## `audit_logs`

Stores actions performed by users.

Important fields:

```text
id
user_id
action
entity_type
entity_id
created_at
```

Relationship:

```text
users.id
   ↑
   │ FOREIGN KEY
   │
audit_logs.user_id
```

---

# Topics Practiced

## Basic Queries

- `SELECT`
- `WHERE`
- `ORDER BY`
- `LIMIT`

## Schema Design

- `CREATE TABLE`
- `PRIMARY KEY`
- `FOREIGN KEY`
- `NOT NULL`
- `UNIQUE`
- `DEFAULT`

## Intermediate Queries

- `JOIN`
- `LEFT JOIN`
- `GROUP BY`
- `COUNT`
- `IS NULL`
- aliases with `AS`

## Database Integration

- PostgreSQL schema initialization
- MySQL schema creation
- foreign-key relationships
- parameterized query preparation
- repository-oriented database design

---

# 10 Query Exercises

`task-system-queries.sql` contains exercises such as:

1. Select all users
2. Select unfinished tasks
3. Select the latest tasks
4. Join tasks with their owners
5. Count tasks by status
6. Count tasks per user
7. View task status history
8. Find specific audit actions
9. Find tasks without status history
10. Count audit actions by type

---

# Users and Tasks JOIN

`users-tasks-join.sql` practices an `INNER JOIN` between:

```text
users
and
tasks
```

The relationship is:

```text
users.id
   ↑
   │
tasks.user_id
```

Example:

```sql
SELECT
    u.id AS user_id,
    u.name AS user_name,
    t.id AS task_id,
    t.title,
    t.status
FROM users AS u
JOIN tasks AS t
    ON t.user_id = u.id
ORDER BY u.id, t.id;
```

The key condition is:

```sql
ON t.user_id = u.id
```

Conceptually:

```text
users
   │
   │ users.id = tasks.user_id
   ▼
tasks
```

Example data:

```text
users

1 Alice
2 Bob
```

```text
tasks

1 Learn Go       user_id=1
2 Practice SQL   user_id=1
3 Learn Docker   user_id=2
```

After the JOIN:

```text
Alice | Learn Go
Alice | Practice SQL
Bob   | Learn Docker
```

The main idea is:

> A JOIN combines related rows from different tables using their relationship columns.

---

# SQLite Practice

Basic SQL exercises can be tested locally with SQLite.

Create the practice database:

```bash
sqlite3 practice.db < 02-sql/setup.sql
```

Run the basic query exercises:

```bash
sqlite3 -header -column practice.db \
  < 02-sql/select-basics.sql
```

The SQLite database is intended only for local practice and should not be committed to Git.

---

# PostgreSQL Practice with Docker Compose

The task system schema is used by the PostgreSQL container defined in:

```text
03-docker/compose.yaml
```

From the repository root, start PostgreSQL:

```bash
docker compose \
  -f 03-docker/compose.yaml \
  up -d
```

Check the container:

```bash
docker compose \
  -f 03-docker/compose.yaml \
  ps
```

Expected status:

```text
Up ... (healthy)
```

Check whether PostgreSQL is ready:

```bash
docker compose \
  -f 03-docker/compose.yaml \
  exec db \
  pg_isready -U app -d taskdb
```

Open PostgreSQL:

```bash
docker compose \
  -f 03-docker/compose.yaml \
  exec db \
  psql -U app -d taskdb
```

List tables:

```sql
\dt
```

Expected tables:

```text
users
tasks
status_history
audit_logs
```

Exit PostgreSQL:

```text
\q
```

---

# Run the PostgreSQL Query Exercises

Run all 10 queries from the repository root:

```bash
docker compose \
  -f 03-docker/compose.yaml \
  exec -T db \
  psql -U app -d taskdb \
  < 02-sql/task-system-queries.sql
```

Run the `users` and `tasks` JOIN:

```bash
docker compose \
  -f 03-docker/compose.yaml \
  exec -T db \
  psql -U app -d taskdb \
  < 02-sql/users-tasks-join.sql
```

Execute an individual query:

```bash
docker compose \
  -f 03-docker/compose.yaml \
  exec db \
  psql -U app -d taskdb \
  -c 'SELECT * FROM tasks;'
```

---

# Stop PostgreSQL

Stop and remove the Compose container and network while keeping the database volume:

```bash
docker compose \
  -f 03-docker/compose.yaml \
  down
```

The persistent PostgreSQL data remains available.

Start it again later with:

```bash
docker compose \
  -f 03-docker/compose.yaml \
  up -d
```

If only stopping the container is desired:

```bash
docker compose \
  -f 03-docker/compose.yaml \
  stop
```

Start it again with:

```bash
docker compose \
  -f 03-docker/compose.yaml \
  start
```

---

# Reset the PostgreSQL Database

The PostgreSQL container uses a persistent Docker volume.

If the schema or seed files change, the initialization scripts are not automatically executed again while the existing volume remains.

To completely recreate the learning database:

```bash
docker compose \
  -f 03-docker/compose.yaml \
  down -v

docker compose \
  -f 03-docker/compose.yaml \
  up -d
```

Warning:

```text
docker compose down -v
```

deletes the PostgreSQL volume and all stored database data.

This is acceptable for the current learning environment because the database can be recreated from:

```text
task-system-schema.sql
task-system-seed.sql
```

For normal shutdown, prefer:

```bash
docker compose \
  -f 03-docker/compose.yaml \
  down
```

without `-v`.

---

# MySQL Task API Schema

`mysql-task-api-schema.sql` defines the table used by the Go MySQL repository exercise.

Schema:

```sql
CREATE TABLE IF NOT EXISTS api_tasks (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(255) NOT NULL,
    done BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

This schema is used by:

```text
01-go/06-database/02-mysql-repository
```

The architecture is:

```text
HTTP Handler
↓
TaskRepository
↓
MySQLTaskRepository
↓
database/sql
↓
MySQL
```

Create the table in the local MySQL practice container:

```bash
docker exec -i mysql-practice \
  mysql \
  -uapp \
  -papp \
  taskdb \
  < 02-sql/mysql-task-api-schema.sql
```

Check the table:

```bash
docker exec mysql-practice \
  mysql \
  -uapp \
  -papp \
  taskdb \
  -e "SHOW TABLES;"
```

Inspect the schema:

```bash
docker exec mysql-practice \
  mysql \
  -uapp \
  -papp \
  taskdb \
  -e "DESCRIBE api_tasks;"
```

---

# Parameterized Queries

The MySQL repository uses parameterized SQL.

Example:

```sql
INSERT INTO api_tasks (title, done)
VALUES (?, ?)
```

Go passes values separately:

```go
db.ExecContext(
	ctx,
	query,
	input.Title,
	input.Done,
)
```

For a single task:

```sql
SELECT id, title, done
FROM api_tasks
WHERE id = ?
```

The important rule is:

```text
SQL structure
+
query arguments
```

instead of:

```text
SQL
+
user input concatenated into the SQL string
```

This reduces SQL injection risk.

Do not build SQL using:

```text
string concatenation
fmt.Sprintf with user input
```

Use placeholders instead.

---

# Key Concepts

## Primary Key

A primary key uniquely identifies one row.

Example:

```sql
id INTEGER PRIMARY KEY
```

or in MySQL:

```sql
id BIGINT PRIMARY KEY AUTO_INCREMENT
```

---

## Foreign Key

A foreign key creates a relationship between tables.

Example:

```sql
FOREIGN KEY (task_id)
REFERENCES tasks(id)
```

The referenced column must normally be:

```text
PRIMARY KEY
or
UNIQUE
```

Correct relationship:

```text
tasks.id
   ↑
   │ FOREIGN KEY
   │
status_history.task_id
```

---

## JOIN

A JOIN combines related rows from different tables.

Example:

```sql
FROM users AS u
JOIN tasks AS t
    ON t.user_id = u.id
```

Relationship:

```text
users.id
   ↑
   │
tasks.user_id
```

---

## INNER JOIN vs LEFT JOIN

`JOIN` / `INNER JOIN`:

```text
only rows with a match
```

Example:

```text
users with tasks
```

`LEFT JOIN`:

```text
all rows from the left table
+
matching rows from the right table
```

Example:

```text
all users
including users with zero tasks
```

---

# Current Relationships

```text
users
  │
  ├── tasks
  │     │
  │     └── status_history
  │
  └── audit_logs
```

More explicitly:

```text
users.id
   │
   └── tasks.user_id

tasks.id
   │
   └── status_history.task_id

users.id
   │
   └── audit_logs.user_id
```

---

# Key Takeaways

```text
Primary Key
→ uniquely identifies a row

Foreign Key
→ references another table

JOIN
→ combines related rows

GROUP BY
→ groups rows for aggregation

COUNT
→ counts rows

Parameterized Query
→ separates SQL structure from user input
```

The most important relationship currently practiced is:

```text
users.id
   ↑
   │
tasks.user_id
```

and the corresponding JOIN:

```sql
JOIN tasks AS t
    ON t.user_id = u.id
```