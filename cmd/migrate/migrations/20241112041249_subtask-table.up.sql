CREATE TABLE IF NOT EXISTS subtasks (
    id SERIAL PRIMARY KEY,
    taskId INT NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
    statusId INT NOT NULL REFERENCES statuses(id) ON DELETE
    SET NULL,
        title TEXT NOT NULL,
        createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        deletedAt TIMESTAMP,
        deletedBy INT REFERENCES users(id) ON DELETE
    SET NULL
);