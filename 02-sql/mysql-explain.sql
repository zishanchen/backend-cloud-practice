-- MySQL index and EXPLAIN practice.

EXPLAIN
SELECT
    id,
    title,
    done
FROM api_tasks
WHERE title = 'Learn Go';

CREATE INDEX idx_api_tasks_title
    ON api_tasks (title);

ANALYZE TABLE api_tasks;

EXPLAIN
SELECT
    id,
    title,
    done
FROM api_tasks
WHERE title = 'Learn Go';

SHOW INDEX
FROM api_tasks;