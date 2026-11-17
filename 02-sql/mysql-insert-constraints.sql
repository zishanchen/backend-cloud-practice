-- MySQL INSERT / Constraints / DEFAULT practice

DROP TABLE IF EXISTS practice_tasks;
DROP TABLE IF EXISTS practice_users;

CREATE TABLE practice_users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE practice_tasks (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(255) NOT NULL,
    status VARCHAR(32) NOT NULL DEFAULT 'todo',
    priority INT NOT NULL DEFAULT 3,
    user_id BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_practice_tasks_user
        FOREIGN KEY (user_id)
        REFERENCES practice_users(id),

    CONSTRAINT chk_practice_tasks_priority
        CHECK (priority BETWEEN 1 AND 5)
);

-- INSERT users

INSERT INTO practice_users (name, email)
VALUES
    ('Alice', 'alice@example.com'),
    ('Bob', 'bob@example.com'),
    ('Charlie', 'charlie@example.com');


-- INSERT tasks

-- status and priority are explicitly provided
INSERT INTO practice_tasks (
    title,
    status,
    priority,
    user_id
)
VALUES
    ('Learn Go', 'done', 2, 1),
    ('Practice SQL', 'in_progress', 1, 1),
    ('Learn Docker', 'todo', 3, 2),
    ('Review HTTP', 'done', 2, 2);

-- status and priority are omitted
-- DEFAULT values will be used
INSERT INTO practice_tasks (
    title,
    user_id
)
VALUES
    ('Learn Kubernetes', 1),
    ('Practice GROUP BY', 2),
    ('Read MySQL docs', 3);

-- Only status is explicitly provided.
-- priority still uses its DEFAULT value.
INSERT INTO practice_tasks (
    title,
    status,
    user_id
)
VALUES
    ('Review JOIN', 'in_progress', 3);


-- Inspect inserted data and DEFAULT values
SELECT
    id,
    title,
    status,
    priority,
    user_id,
    created_at
FROM practice_tasks
ORDER BY id;
