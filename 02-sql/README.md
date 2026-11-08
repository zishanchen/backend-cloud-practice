# SQL Practice

This folder contains SQL exercises, schema design practice, sample data, and query exercises used with SQLite and PostgreSQL.

## Files

- `setup.sql` — creates sample data for basic SQL practice
- `select-basics.sql` — practices basic `SELECT`, `WHERE`, `ORDER BY`, and `LIMIT`
- `task-system-schema.sql` — defines the task management database schema
- `task-system-seed.sql` — inserts sample data into the task management tables
- `task-system-queries.sql` — contains 10 SQL queries for practicing filtering, joins, grouping, and aggregation

## Task System Schema

The task system contains four tables:

```text
users
tasks
status_history
audit_logs
```

### `users`

Stores application users.

Important fields:

```text
id
name
email
```

### `tasks`

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

### `status_history`

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

### `audit_logs`

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

## Topics Practiced

### Basic Queries

- `SELECT`
- `WHERE`
- `ORDER BY`
- `LIMIT`

### Schema Design

- `CREATE TABLE`
- `PRIMARY KEY`
- `FOREIGN KEY`
- `NOT NULL`
- `UNIQUE`
- `DEFAULT`

### Intermediate Queries

- `JOIN`
- `LEFT JOIN`
- `GROUP BY`
- `COUNT`
- `IS NULL`
- aliases with `AS`

## 10 Query Exercises

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

## SQLite Practice

Basic SQL exercises can be tested locally with SQLite.

Example:

```bash
sqlite3 practice.db < 02-sql/setup.sql
```

Run the basic query exercises:

```bash
sqlite3 -header -column practice.db \
  < 02-sql/select-basics.sql
```

## PostgreSQL Practice with Docker Compose

The task system schema is also used by the PostgreSQL container defined in:

```text
03-docker/compose.yaml
```

Start PostgreSQL:

```bash
cd 03-docker
docker compose up -d
```

Check the container:

```bash
docker compose ps
```

Open PostgreSQL:

```bash
docker compose exec db \
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

## Run the SQL Query Exercises

From the repository root:

```bash
docker exec -i task-postgres \
  psql -U app -d taskdb \
  < 02-sql/task-system-queries.sql
```

Or execute individual queries:

```bash
docker compose exec db \
  psql -U app -d taskdb \
  -c 'SELECT * FROM tasks;'
```

## Reset the PostgreSQL Database

The PostgreSQL container uses a persistent Docker volume.

If the schema or seed files are changed, the initialization scripts will not automatically run again while the existing volume remains.

To recreate the database:

```bash
cd 03-docker

docker compose down -v
docker compose up -d
```

Warning:

```text
docker compose down -v
```

deletes the PostgreSQL data volume and all stored database data.

This is acceptable for the current learning environment because the database can be recreated from:

```text
task-system-schema.sql
task-system-seed.sql
```

## Key Concepts

### Primary Key

A primary key uniquely identifies one row.

Example:

```sql
id INTEGER PRIMARY KEY
```

### Foreign Key

A foreign key creates a relationship between tables.

Example:

```sql
FOREIGN KEY (task_id)
REFERENCES tasks(id)
```

The referenced column must normally be a:

```text
PRIMARY KEY
or
UNIQUE
```

### Current Relationships

```text
users
  │
  ├── tasks
  │     │
  │     └── status_history
  │
  └── audit_logs
```

The main idea is:

```text
Primary Key
→ identifies a row

Foreign Key
→ references another row
```