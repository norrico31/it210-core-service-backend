CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    description TEXT,
    statusId INT REFERENCES statuses(id) ON DELETE
    SET NULL,
        userId INT REFERENCES users(id) ON DELETE
    SET NULL,
        projectId INT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
        priorityId INT NOT NULL REFERENCES priorities(id) ON DELETE CASCADE,
        workspaceId INT REFERENCES workspaces(id) ON DELETE CASCADE,
        taskOrder INT,
        taskProgress INT,
        createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        deletedAt TIMESTAMP,
        deletedBy INT REFERENCES users(id) ON DELETE
    SET NULL
);