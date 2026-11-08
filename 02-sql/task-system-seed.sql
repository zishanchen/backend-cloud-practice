INSERT INTO users (id, name, email) VALUES
(1, 'Alice', 'alice@example.com'),
(2, 'Bob', 'bob@example.com'),
(3, 'Charlie', 'charlie@example.com');

INSERT INTO tasks (
    id,
    title,
    status,
    user_id
) VALUES
(1, 'Learn Go', 'done', 1),
(2, 'Practice SQL', 'in_progress', 1),
(3, 'Learn Docker', 'todo', 2),
(4, 'Review HTTP', 'done', 2),
(5, 'Learn Compose', 'todo', 3);

INSERT INTO status_history (
    id,
    task_id,
    old_status,
    new_status
) VALUES
(1, 1, 'todo', 'in_progress'),
(2, 1, 'in_progress', 'done'),
(3, 2, 'todo', 'in_progress'),
(4, 4, 'todo', 'done');

INSERT INTO audit_logs (
    id,
    user_id,
    action,
    entity_type,
    entity_id
) VALUES
(1, 1, 'CREATE_TASK', 'task', 1),
(2, 1, 'UPDATE_TASK', 'task', 1),
(3, 1, 'CREATE_TASK', 'task', 2),
(4, 2, 'CREATE_TASK', 'task', 3),
(5, 2, 'UPDATE_TASK', 'task', 4);