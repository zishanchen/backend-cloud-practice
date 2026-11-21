package main

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func openIntegrationDB(t *testing.T) *sql.DB {
	t.Helper()

	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		t.Skip("MYSQL_DSN is not set; skipping MySQL integration test")
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("open database: %v", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		t.Fatalf("ping database: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func TestMySQLTaskRepositoryUpdateDoneTransaction(t *testing.T) {
	db := openIntegrationDB(t)
	ctx := context.Background()

	// Create a task specifically for this integration test.
	result, err := db.ExecContext(
		ctx,
		`
			INSERT INTO api_tasks (
				title,
				done
			)
			VALUES (?, ?)
		`,
		"transaction integration test",
		false,
	)
	if err != nil {
		t.Fatalf("insert test task: %v", err)
	}

	taskID, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("get task id: %v", err)
	}

	t.Cleanup(func() {
		_, _ = db.Exec(
			`
				DELETE FROM api_tasks
				WHERE id = ?
			`,
			taskID,
		)
	})

	repo := NewMySQLTaskRepository(db)

	// --------------------------------------------------
	// 1. Successful transaction
	// --------------------------------------------------

	task, err := repo.UpdateDone(
		ctx,
		taskID,
		true,
	)
	if err != nil {
		t.Fatalf("update done: %v", err)
	}

	if !task.Done {
		t.Fatal("expected returned task to be done")
	}

	// Verify that api_tasks was updated.
	var databaseDone bool

	err = db.QueryRowContext(
		ctx,
		`
			SELECT done
			FROM api_tasks
			WHERE id = ?
		`,
		taskID,
	).Scan(&databaseDone)
	if err != nil {
		t.Fatalf("read updated task: %v", err)
	}

	if !databaseDone {
		t.Fatal("expected database task to be done")
	}

	// Verify that the audit log was created.
	var (
		oldDone bool
		newDone bool
	)

	err = db.QueryRowContext(
		ctx,
		`
			SELECT
				old_done,
				new_done
			FROM api_task_audit_logs
			WHERE task_id = ?
			ORDER BY id DESC
			LIMIT 1
		`,
		taskID,
	).Scan(
		&oldDone,
		&newDone,
	)
	if err != nil {
		t.Fatalf("read audit log: %v", err)
	}

	if oldDone {
		t.Fatal("expected old_done=false")
	}

	if !newDone {
		t.Fatal("expected new_done=true")
	}

	// --------------------------------------------------
	// 2. Conflict test
	// --------------------------------------------------

	_, err = repo.UpdateDone(
		ctx,
		taskID,
		true,
	)

	if !errors.Is(err, ErrTaskConflict) {
		t.Fatalf(
			"expected ErrTaskConflict, got %v",
			err,
		)
	}

	// --------------------------------------------------
	// 3. Confirm that the failed/conflicting operation
	// --------------------------------------------------

	var auditCount int

	err = db.QueryRowContext(
		ctx,
		`
			SELECT COUNT(*)
			FROM api_task_audit_logs
			WHERE task_id = ?
		`,
		taskID,
	).Scan(&auditCount)
	if err != nil {
		t.Fatalf("count audit logs: %v", err)
	}

	if auditCount != 1 {
		t.Fatalf(
			"expected 1 audit log, got %d",
			auditCount,
		)
	}
}
