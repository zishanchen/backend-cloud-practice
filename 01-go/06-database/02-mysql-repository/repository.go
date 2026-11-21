package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

var (
	ErrTaskNotFound = errors.New("task not found")
	ErrTaskConflict = errors.New("task state conflict")
)

type TaskRepository interface {
	Create(
		ctx context.Context,
		input CreateTaskRequest,
	) (Task, error)

	Get(
		ctx context.Context,
		id int64,
	) (Task, error)

	List(
		ctx context.Context,
	) ([]Task, error)

	UpdateDone(
		ctx context.Context,
		id int64,
		done bool,
	) (Task, error)
}

type MySQLTaskRepository struct {
	db *sql.DB
}

func NewMySQLTaskRepository(
	db *sql.DB,
) *MySQLTaskRepository {
	return &MySQLTaskRepository{
		db: db,
	}
}

func (r *MySQLTaskRepository) Create(
	ctx context.Context,
	input CreateTaskRequest,
) (Task, error) {
	const query = `
		INSERT INTO api_tasks (title, done)
		VALUES (?, ?)
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		input.Title,
		input.Done,
	)
	if err != nil {
		log.Printf("MySQL INSERT error: %v", err)

		return Task{}, fmt.Errorf(
			"create task: %w",
			err,
		)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Task{}, fmt.Errorf(
			"get inserted task id: %w",
			err,
		)
	}

	return Task{
		ID:    id,
		Title: input.Title,
		Done:  input.Done,
	}, nil
}

func (r *MySQLTaskRepository) Get(
	ctx context.Context,
	id int64,
) (Task, error) {
	const query = `
		SELECT id, title, done
		FROM api_tasks
		WHERE id = ?
	`

	var task Task

	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&task.ID,
		&task.Title,
		&task.Done,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return Task{}, ErrTaskNotFound
	}

	if err != nil {
		return Task{}, fmt.Errorf(
			"get task: %w",
			err,
		)
	}

	return task, nil
}

func (r *MySQLTaskRepository) List(
	ctx context.Context,
) ([]Task, error) {
	const query = `
		SELECT id, title, done
		FROM api_tasks
		ORDER BY id
	`

	rows, err := r.db.QueryContext(
		ctx,
		query,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"list tasks: %w",
			err,
		)
	}
	defer rows.Close()

	tasks := []Task{}

	for rows.Next() {
		var task Task

		if err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Done,
		); err != nil {
			return nil, fmt.Errorf(
				"scan task: %w",
				err,
			)
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"iterate tasks: %w",
			err,
		)
	}

	return tasks, nil
}

func (r *MySQLTaskRepository) UpdateDone(
	ctx context.Context,
	id int64,
	done bool,
) (Task, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return Task{}, fmt.Errorf(
			"begin update transaction: %w",
			err,
		)
	}

	defer tx.Rollback()

	var task Task

	err = tx.QueryRowContext(
		ctx,
		`
			SELECT
				id,
				title,
				done
			FROM api_tasks
			WHERE id = ?
			FOR UPDATE
		`,
		id,
	).Scan(
		&task.ID,
		&task.Title,
		&task.Done,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return Task{}, ErrTaskNotFound
	}

	if err != nil {
		return Task{}, fmt.Errorf(
			"read task for update: %w",
			err,
		)
	}

	if task.Done == done {
		return Task{}, ErrTaskConflict
	}

	oldDone := task.Done

	_, err = tx.ExecContext(
		ctx,
		`
			UPDATE api_tasks
			SET done = ?
			WHERE id = ?
		`,
		done,
		id,
	)
	if err != nil {
		return Task{}, fmt.Errorf(
			"update task state: %w",
			err,
		)
	}

	_, err = tx.ExecContext(
		ctx,
		`
			INSERT INTO api_task_audit_logs (
				task_id,
				old_done,
				new_done
			)
			VALUES (?, ?, ?)
		`,
		id,
		oldDone,
		done,
	)
	if err != nil {
		return Task{}, fmt.Errorf(
			"insert task audit log: %w",
			err,
		)
	}

	if err := tx.Commit(); err != nil {
		return Task{}, fmt.Errorf(
			"commit update transaction: %w",
			err,
		)
	}

	task.Done = done

	return task, nil
}
