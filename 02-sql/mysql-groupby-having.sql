-- MySQL GROUP BY / HAVING practice

-- 1. Count tasks by status

SELECT
    status,
    COUNT(*) AS task_count
FROM practice_tasks
GROUP BY status
ORDER BY task_count DESC;

-- 2. Find users who have at least 2 tasks

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


-- 3. Find users who have at least 1 completed task

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
) >= 1
ORDER BY completed_tasks DESC;