# SQL Practice

This folder contains SQL exercises, schema design practice, sample data, JOIN exercises, aggregation queries, constraints, and database schemas used with SQLite, PostgreSQL, and MySQL.

## Files

- `setup.sql` — creates sample data for basic SQLite SQL practice
- `select-basics.sql` — practices basic `SELECT`, `WHERE`, `ORDER BY`, and `LIMIT`
- `task-system-schema.sql` — defines the PostgreSQL task management database schema
- `task-system-seed.sql` — inserts sample data into the task management tables
- `task-system-queries.sql` — contains 10 SQL queries for filtering, joins, grouping, and aggregation
- `users-tasks-join.sql` — practices joining `users` and `tasks`
- `mysql-task-api-schema.sql` — defines the MySQL table used by the Go repository practice
- `mysql-transaction-schema.sql` — defines tables used for Go transaction and rollback practice
- `mysql-insert-constraints.sql` — practices `INSERT`, `DEFAULT`, and database constraints
- `mysql-groupby-having.sql` — practices `GROUP BY`, aggregation, and `HAVING`

---

# Topics Practiced

## Basic Queries

- `SELECT`
- `WHERE`
- `ORDER BY`
- `LIMIT`

## Data Manipulation

- `INSERT INTO`
- multi-row `INSERT`
- omitted columns
- `DEFAULT` values

## Schema Design

- `CREATE TABLE`
- `PRIMARY KEY`
- `AUTO_INCREMENT`
- `FOREIGN KEY`
- `NOT NULL`
- `UNIQUE`
- `CHECK`
- `DEFAULT`

## Intermediate Queries

- `JOIN`
- `INNER JOIN`
- `LEFT JOIN`
- aliases with `AS`
- `IS NULL`

## Aggregation

- `GROUP BY`
- `HAVING`
- `COUNT`
- `SUM`
- `CASE`

## Database Integration

- SQLite local practice
- PostgreSQL with Docker Compose
- MySQL with Docker
- parameterized queries
- SQL injection prevention
- repository-oriented database design
- transaction-related schemas

---

# PostgreSQL Task System

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

Each task belongs to one user.

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

Each status history record belongs to one task.

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

Relationship:

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

A JOIN combines related rows from different tables using their relationship columns.

---

# MySQL INSERT and Constraints

`mysql-insert-constraints.sql` practices inserting data and enforcing database-level rules.

Example user table:

```sql
CREATE TABLE practice_users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

Example task table:

```sql
CREATE TABLE practice_tasks (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(255) NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'todo',
    priority INT NOT NULL DEFAULT 3,
    user_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_practice_tasks_user
        FOREIGN KEY (user_id)
        REFERENCES practice_users(id),

    CONSTRAINT chk_practice_tasks_priority
        CHECK (priority BETWEEN 1 AND 5)
);
```

---

# INSERT

Basic example:

```sql
INSERT INTO practice_users (
    name,
    email
)
VALUES (
    'Alice',
    'alice@example.com'
);
```

Multiple rows can be inserted together:

```sql
INSERT INTO practice_users (
    name,
    email
)
VALUES
    ('Alice', 'alice@example.com'),
    ('Bob', 'bob@example.com'),
    ('Charlie', 'charlie@example.com');
```

---

# DEFAULT Values

Suppose:

```sql
status VARCHAR(32)
    NOT NULL
    DEFAULT 'todo',

priority INT
    NOT NULL
    DEFAULT 3
```

Then:

```sql
INSERT INTO practice_tasks (
    title,
    user_id
)
VALUES (
    'Learn Kubernetes',
    1
);
```

does not specify:

```text
status
priority
```

so MySQL automatically uses:

```text
status   = todo
priority = 3
```

The idea is:

```text
column omitted
↓
DEFAULT exists
↓
database inserts default value
```

---

# Database Constraints

Constraints protect data integrity at the database level.

## `PRIMARY KEY`

Uniquely identifies one row.

```sql
id BIGINT PRIMARY KEY AUTO_INCREMENT
```

## `NOT NULL`

Prevents a column from containing `NULL`.

```sql
title VARCHAR(255) NOT NULL
```

## `UNIQUE`

Prevents duplicate values.

```sql
email VARCHAR(255) NOT NULL UNIQUE
```

For example, inserting the same email twice should fail.

## `FOREIGN KEY`

Ensures that a referenced row exists.

```sql
FOREIGN KEY (user_id)
REFERENCES practice_users(id)
```

Relationship:

```text
practice_users.id
        ↑
        │ FOREIGN KEY
        │
practice_tasks.user_id
```

If:

```text
user_id = 999
```

but user `999` does not exist, the database rejects the row.

## `CHECK`

Ensures a value satisfies a condition.

```sql
CHECK (priority BETWEEN 1 AND 5)
```

Valid:

```text
1
2
3
4
5
```

Invalid:

```text
0
10
-1
```

## `DEFAULT`

Provides a value when a column is omitted.

```sql
status VARCHAR(32)
    NOT NULL
    DEFAULT 'todo'
```

---

# Constraint Naming: `fk_` and `chk_`

Names such as:

```text
fk_practice_tasks_user
chk_practice_tasks_priority
```

are user-defined constraint names.

`fk_` is a common abbreviation for:

```text
Foreign Key
```

Example:

```sql
CONSTRAINT fk_practice_tasks_user
    FOREIGN KEY (user_id)
    REFERENCES practice_users(id)
```

`chk_` is a common abbreviation for:

```text
CHECK
```

Example:

```sql
CONSTRAINT chk_practice_tasks_priority
    CHECK (priority BETWEEN 1 AND 5)
```

These prefixes are naming conventions, not SQL keywords.

A useful convention is:

```text
pk_   → PRIMARY KEY
fk_   → FOREIGN KEY
uq_   → UNIQUE
chk_  → CHECK
```

---

# GROUP BY

`GROUP BY` combines rows with the same grouping value.

Example:

```sql
SELECT
    status,
    COUNT(*) AS task_count
FROM practice_tasks
GROUP BY status;
```

Conceptually:

```text
all tasks
↓
group by status
↓
todo group
done group
in_progress group
↓
COUNT each group
```

Possible result:

```text
status        task_count
------------  ----------
done          2
in_progress   2
todo          4
```

---

# HAVING

`HAVING` filters groups after aggregation.

Example:

```sql
SELECT
    user_id,
    COUNT(*) AS task_count
FROM practice_tasks
GROUP BY user_id
HAVING COUNT(*) >= 2;
```

This means:

> Group tasks by user, count each group, and keep only users with at least two tasks.

---

# WHERE vs HAVING

This distinction is important.

## `WHERE`

Filters individual rows before grouping.

```sql
SELECT
    status,
    COUNT(*)
FROM practice_tasks
WHERE priority <= 3
GROUP BY status;
```

Flow:

```text
rows
↓
WHERE
↓
GROUP BY
↓
aggregation
```

## `HAVING`

Filters groups after grouping and aggregation.

```sql
SELECT
    user_id,
    COUNT(*) AS task_count
FROM practice_tasks
GROUP BY user_id
HAVING COUNT(*) >= 2;
```

Flow:

```text
rows
↓
GROUP BY
↓
COUNT
↓
HAVING
```

Remember:

```text
WHERE
→ filters rows

HAVING
→ filters groups
```

---

# GROUP BY / HAVING Exercises

`mysql-groupby-having.sql` contains three main exercises.

## 1. Count Tasks by Status

```sql
SELECT
    status,
    COUNT(*) AS task_count
FROM practice_tasks
GROUP BY status
ORDER BY task_count DESC;
```

Practices:

```text
GROUP BY
COUNT
ORDER BY
```

---

## 2. Users with at Least Two Tasks

```sql
SELECT
    u.id AS user_id,
    u.name AS user_name,
    COUNT(t.id) AS task_count
FROM practice_users AS u
JOIN practice_tasks AS t
    ON t.user_id = u.id
GROUP BY
    u.id,
    u.name
HAVING COUNT(t.id) >= 2
ORDER BY task_count DESC;
```

Practices:

```text
JOIN
GROUP BY
COUNT
HAVING
```

---

## 3. Users with at Least One Completed Task

```sql
SELECT
    u.id AS user_id,
    u.name AS user_name,
    COUNT(t.id) AS total_tasks,
    SUM(
        CASE
            WHEN t.status = 'done' THEN 1
            ELSE 0
        END
    ) AS completed_tasks
FROM practice_users AS u
JOIN practice_tasks AS t
    ON t.user_id = u.id
GROUP BY
    u.id,
    u.name
HAVING SUM(
    CASE
        WHEN t.status = 'done' THEN 1
        ELSE 0
    END
) >= 1;
```

The expression:

```sql
CASE
    WHEN t.status = 'done' THEN 1
    ELSE 0
END
```

converts:

```text
done
→ 1

not done
→ 0
```

Then:

```sql
SUM(...)
```

counts completed tasks.

---

# MySQL Task API Schema

`mysql-task-api-schema.sql` defines the table used by:

```text
01-go/06-database/02-mysql-repository
```

Schema:

```sql
CREATE TABLE IF NOT EXISTS api_tasks (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(255) NOT NULL,
    done BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

Architecture:

```text
HTTP Handler
↓
Service
↓
TaskRepository
↓
MySQLTaskRepository
↓
database/sql
↓
MySQL
```

---

# Parameterized Queries

The Go MySQL repository uses parameterized SQL.

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

For lookup:

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
user input concatenated into SQL text
```

This reduces SQL injection risk.

Do not build SQL using user input with:

```text
string concatenation
fmt.Sprintf
```

Use placeholders instead.

---

# MySQL Transaction Schema

`mysql-transaction-schema.sql` is used by the Go database transaction exercise.

The main relationship is:

```text
tx_tasks
   │
   └── tx_audit_logs
```

The application performs:

```text
UPDATE task status
+
INSERT audit log
```

inside the same transaction.

Expected behavior:

```text
both succeed
→ COMMIT
```

or:

```text
one fails
→ ROLLBACK
```

This keeps related database changes atomic.

---

# SQLite Practice

Create the local SQLite practice database:

```bash
sqlite3 practice.db < 02-sql/setup.sql
```

Run basic queries:

```bash
sqlite3 -header -column practice.db \
  < 02-sql/select-basics.sql
```

`practice.db` is local practice data and should not be committed to Git.

---

# PostgreSQL with Docker Compose

Start PostgreSQL:

```bash
docker compose \
  -f 03-docker/compose.yaml \
  up -d
```

Check:

```bash
docker compose \
  -f 03-docker/compose.yaml \
  ps
```

Check readiness:

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

Expected task-system tables:

```text
users
tasks
status_history
audit_logs
```

---

# Run PostgreSQL Exercises

Run the 10 query exercises:

```bash
docker compose \
  -f 03-docker/compose.yaml \
  exec -T db \
  psql -U app -d taskdb \
  < 02-sql/task-system-queries.sql
```

Run the users/tasks JOIN:

```bash
docker compose \
  -f 03-docker/compose.yaml \
  exec -T db \
  psql -U app -d taskdb \
  < 02-sql/users-tasks-join.sql
```

---

# Stop PostgreSQL

Stop and remove the Compose container/network while keeping the volume:

```bash
docker compose \
  -f 03-docker/compose.yaml \
  down
```

Start again:

```bash
docker compose \
  -f 03-docker/compose.yaml \
  up -d
```

Do not normally use:

```bash
docker compose \
  -f 03-docker/compose.yaml \
  down -v
```

because `-v` deletes the PostgreSQL data volume.

---

# MySQL Practice

Start the local MySQL container:

```bash
docker start mysql-practice
```

Check:

```bash
docker ps \
  --filter name=mysql-practice
```

Run the INSERT / constraint setup:

```bash
docker exec -i mysql-practice \
  mysql \
  -uapp \
  -papp \
  taskdb \
  < 02-sql/mysql-insert-constraints.sql
```

Run the aggregation exercises:

```bash
docker exec -i mysql-practice \
  mysql \
  -uapp \
  -papp \
  taskdb \
  < 02-sql/mysql-groupby-having.sql
```

Inspect the practice table:

```bash
docker exec mysql-practice \
  mysql \
  -uapp \
  -papp \
  taskdb \
  -e "SHOW CREATE TABLE practice_tasks\G"
```

---

# Current Relationships

PostgreSQL task system:

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

MySQL constraint practice:

```text
practice_users.id
       │
       └── practice_tasks.user_id
```

MySQL transaction practice:

```text
tx_tasks.id
   │
   └── tx_audit_logs.task_id
```

---

# Key Takeaways

```text
SELECT
→ reads rows

INSERT
→ creates rows

PRIMARY KEY
→ uniquely identifies a row

FOREIGN KEY
→ enforces a relationship

NOT NULL
→ prevents NULL values

UNIQUE
→ prevents duplicate values

CHECK
→ validates column values

DEFAULT
→ supplies an omitted value

JOIN
→ combines related rows

GROUP BY
→ groups rows

COUNT / SUM
→ aggregate each group

WHERE
→ filters rows before grouping

HAVING
→ filters groups after grouping

Parameterized Query
→ separates SQL structure from user input

Transaction
→ groups related database changes atomically
```

The key relationships practiced are:

```text
users.id
   ↑
   │
tasks.user_id
```

and:

```text
practice_users.id
        ↑
        │ FOREIGN KEY
        │
practice_tasks.user_id
```

The key aggregation distinction is:

```text
WHERE
→ before GROUP BY

HAVING
→ after GROUP BY
```