-- INNER JOIN users and tasks
-- Only users that have tasks appear in the result

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

