CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    subTask TEXT [],
    description TEXT,
    statusId INT NOT NULL REFERENCES statuses(id) ON DELETE
    SET NULL,
        userId INT REFERENCES users(id) ON DELETE
    SET NULL,
        projectId INT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
        createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        deletedAt TIMESTAMP
);