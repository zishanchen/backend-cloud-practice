-- Audit table for transactional task status updates.

CREATE TABLE IF NOT EXISTS api_task_audit_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    task_id BIGINT NOT NULL,
    old_done BOOLEAN NOT NULL,
    new_done BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_api_task_audit_task
        FOREIGN KEY (task_id)
        REFERENCES api_tasks(id)
        ON DELETE CASCADE
);