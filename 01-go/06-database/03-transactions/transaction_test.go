package main

import (
	"context"
	"database/sql"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()

	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		t.Skip(
			"MYSQL_DSN is not set; skipping MySQL integration test",
		)
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf(
			"open database: %v",
			err,
		)
	}

	if err := db.Ping(); err != nil {
		db.Close()

		t.Fatalf(
			"ping database: %v",
			err,
		)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return db
}

func createTestTask(
	t *testing.T,
	db *sql.DB,
	status string,
) int64 {
	t.Helper()

	result, err := db.Exec(
		`
			INSERT INTO tx_tasks (
				title,
				status
			)
			VALUES (?, ?)
		`,
		"transaction test task",
		status,
	)
	if err != nil {
		t.Fatalf(
			"create test task: %v",
			err,
		)
	}

	id, err := result.LastInsertId()
	if err != nil {
		t.Fatalf(
			"get task id: %v",
			err,
		)
	}

	t.Cleanup(func() {
		_, _ = db.Exec(
			`
				DELETE FROM tx_tasks
				WHERE id = ?
			`,
			id,
		)
	})

	return id
}

func TestUpdateTaskStatusCommit(
	t *testing.T,
) {
	db := openTestDB(t)
	ctx := context.Background()

	taskID := createTestTask(
		t,
		db,
		"todo",
	)

	err := UpdateTaskStatus(
		ctx,
		db,
		taskID,
		"done",
		false,
	)
	if err != nil {
		t.Fatalf(
			"update task status: %v",
			err,
		)
	}

	var status string

	err = db.QueryRow(
		`
			SELECT status
			FROM tx_tasks
			WHERE id = ?
		`,
		taskID,
	).Scan(&status)
	if err != nil {
		t.Fatalf(
			"read task status: %v",
			err,
		)
	}

	if status != "done" {
		t.Fatalf(
			"expected status done, got %s",
			status,
		)
	}

	var auditCount int

	err = db.QueryRow(
		`
			SELECT COUNT(*)
			FROM tx_audit_logs
			WHERE task_id = ?
		`,
		taskID,
	).Scan(&auditCount)
	if err != nil {
		t.Fatalf(
			"count audit logs: %v",
			err,
		)
	}

	if auditCount != 1 {
		t.Fatalf(
			"expected 1 audit log, got %d",
			auditCount,
		)
	}
}

func TestUpdateTaskStatusRollback(
	t *testing.T,
) {
	db := openTestDB(t)
	ctx := context.Background()

	taskID := createTestTask(
		t,
		db,
		"todo",
	)

	err := UpdateTaskStatus(
		ctx,
		db,
		taskID,
		"done",
		true,
	)

	if err == nil {
		t.Fatal(
			"expected simulated transaction failure",
		)
	}

	var status string

	err = db.QueryRow(
		`
			SELECT status
			FROM tx_tasks
			WHERE id = ?
		`,
		taskID,
	).Scan(&status)
	if err != nil {
		t.Fatalf(
			"read task status: %v",
			err,
		)
	}

	if status != "todo" {
		t.Fatalf(
			"expected rollback to keep status todo, got %s",
			status,
		)
	}

	var auditCount int

	err = db.QueryRow(
		`
			SELECT COUNT(*)
			FROM tx_audit_logs
			WHERE task_id = ?
		`,
		taskID,
	).Scan(&auditCount)
	if err != nil {
		t.Fatalf(
			"count audit logs: %v",
			err,
		)
	}

	if auditCount != 0 {
		t.Fatalf(
			"expected 0 audit logs after rollback, got %d",
			auditCount,
		)
	}
}
