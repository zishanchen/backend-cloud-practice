package main

import (
	"context"
	"database/sql"
	"fmt"
)

func UpdateTaskStatus(
	ctx context.Context,
	db *sql.DB,
	taskID int64,
	newStatus string,
	simulateAuditFailure bool,
) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf(
			"begin transaction: %w",
			err,
		)
	}

	// Safe even if Commit succeeds.
	defer tx.Rollback()

	var oldStatus string

	err = tx.QueryRowContext(
		ctx,
		`
			SELECT status
			FROM tx_tasks
			WHERE id = ?
			FOR UPDATE
		`,
		taskID,
	).Scan(&oldStatus)
	if err != nil {
		return fmt.Errorf(
			"get current task status: %w",
			err,
		)
	}

	_, err = tx.ExecContext(
		ctx,
		`
			UPDATE tx_tasks
			SET status = ?
			WHERE id = ?
		`,
		newStatus,
		taskID,
	)
	if err != nil {
		return fmt.Errorf(
			"update task status: %w",
			err,
		)
	}

	// Used only to demonstrate rollback in the test.
	if simulateAuditFailure {
		return fmt.Errorf(
			"simulated audit log failure",
		)
	}

	_, err = tx.ExecContext(
		ctx,
		`
			INSERT INTO tx_audit_logs (
				task_id,
				old_status,
				new_status
			)
			VALUES (?, ?, ?)
		`,
		taskID,
		oldStatus,
		newStatus,
	)
	if err != nil {
		return fmt.Errorf(
			"insert audit log: %w",
			err,
		)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf(
			"commit transaction: %w",
			err,
		)
	}

	return nil
}
