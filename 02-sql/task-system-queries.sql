--查询所有用户
SELECT *
FROM users;

--查询所有未完成任务
SELECT *
FROM tasks
WHERE status <> 'done';

-- 查询最近创建的3个任务
SELECT id, title, status, created_at
FROM tasks
ORDER BY created_at DESC
LIMIT 3;

-- 查询任务以及负责人
SELECT
    tasks.id,
    tasks.title,
    tasks.status,
    users.name AS owner
FROM tasks
JOIN users 
    ON tasks.user_id = users.id;

-- 按状态统计任务数量
SELECT
    status,
    COUNT(*) AS task_count
FROM tasks
GROUP BY status;

-- 统计每个用户拥有多少任务
SELECT
    users.name,
    COUNT(tasks.id) AS task_count
FROM users
LEFT JOIN tasks
    ON users.id = tasks.user_id
GROUP BY users.id, users.name;

-- 查询Task 1 的状态变化历史
SELECT
    old_status,
    new_status,
    created_at
FROM status_history
WHERE task_id = 1
ORDER BY changed_at;

-- 查询所有UPDATE_TASK操作
SELECT
    users.name,
    audit_logs.action,
    audit_logs.entity_id,
    audit_logs.created_at
FROM audit_logs
JOIN users
    ON audit_logs.user_id = users.id
WHERE audit_logs.action = 'UPDATE_TASK';

-- 查询没有状态历史的任务
SELECT
    tasks.id,
    tasks.title,
FROM tasks
LEFT JOIN status_history
    ON tasks.id = status_history.task_id
WHERE status_history.id IS NULL;

-- 按action统计audit logs数量
SELECT
    action,
    COUNT(*) AS action_count
FROM audit_logs
GROUP BY action
ORDER BY action_count DESC;