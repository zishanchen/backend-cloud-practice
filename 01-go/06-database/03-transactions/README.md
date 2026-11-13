# Go Database Transactions

This exercise demonstrates how to use `database/sql` transactions to keep multiple database changes atomic.

The example performs two related operations:

```text
update task status
+
insert audit log
```

They must either both succeed or both fail.

## Transaction Flow

```text
BeginTx
↓
read old status
↓
update task status
↓
insert audit log
↓
Commit
```

If any operation fails:

```text
BeginTx
↓
update task status
↓
error
↓
Rollback
↓
database remains unchanged
```

## Core Pattern

```go
tx, err := db.BeginTx(ctx, nil)
if err != nil {
	return err
}

defer tx.Rollback()

// Use tx for every operation.

if err := tx.Commit(); err != nil {
	return err
}
```

## Important Rule

Inside a transaction, use:

```go
tx.ExecContext(...)
tx.QueryContext(...)
tx.QueryRowContext(...)
```

Do not mix transactional operations with:

```go
db.ExecContext(...)
db.QueryContext(...)
```

because those operations may execute outside the transaction.

## Atomicity

The transaction contains:

```text
1. UPDATE tx_tasks
2. INSERT tx_audit_logs
```

Successful case:

```text
UPDATE succeeds
↓
audit INSERT succeeds
↓
COMMIT
↓
both changes persist
```

Failure case:

```text
UPDATE succeeds
↓
audit operation fails
↓
ROLLBACK
↓
neither change persists
```

## Successful Transaction Test

Initial state:

```text
task status = todo
audit logs = 0
```

Transaction:

```text
todo
↓
done
+
audit log
↓
COMMIT
```

Expected result:

```text
task status = done
audit logs = 1
```

## Rollback Test

Initial state:

```text
task status = todo
audit logs = 0
```

Transaction:

```text
UPDATE status → done
↓
simulated audit failure
↓
ROLLBACK
```

Expected result:

```text
task status = todo
audit logs = 0
```

The status update disappears because it was never committed.

## Run

Set the MySQL DSN:

```bash
export MYSQL_DSN='app:app@tcp(127.0.0.1:3306)/taskdb?parseTime=true'
```

Run:

```bash
go test -v \
  ./01-go/06-database/03-transactions
```

## Key Concepts

```text
sql.DB
→ database handle / connection pool

sql.Tx
→ one database transaction

BeginTx
→ start transaction

tx.ExecContext
→ execute changes inside transaction

tx.QueryRowContext
→ query inside transaction

Commit
→ persist every transaction change

Rollback
→ discard every uncommitted transaction change
```

## Key Takeaway

```text
Related database operations
↓
same transaction
↓
all succeed
or
all fail
```

For this example:

```text
status update
+
audit log
=
one atomic operation
```