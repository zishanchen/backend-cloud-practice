-- 1. 查询所有 tasks

SELECT * 
FROM tasks;


-- 2. 只查询 title 和 done

SELECT
    title,
    done
FROM tasks;

-- 3. 查询所有未完成的 tasks

SELECT *
FROM tasks
WHERE done = false;

-- 4. 按 id 从大到小查询所有 tasks

SELECT *
FROM tasks
ORDER BY id DESC;

-- 5. 查询前 2 个未完成的 tasks，
--    只显示 id 和 title，
--    按 id 从小到大排列

SELECT
    id,
    title
FROM tasks
WHERE done = false
ORDER BY id ASC
LIMIT 2;