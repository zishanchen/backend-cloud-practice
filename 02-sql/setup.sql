DROP TABLE IF EXISTS tasks;

CREATE TABLE tasks (
    id INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    done BOOLEAN NOT NULL DEFAULT false
);

INSERT INTO tasks (id, title, done) VALUES
    (1, 'Learn GO', true),
    (2, 'Practice SQL', false),
    (3, 'Learn HTTP', true),
    (4, 'Review Linked List', false),
    (5, 'Learn Kubernetes', false);

