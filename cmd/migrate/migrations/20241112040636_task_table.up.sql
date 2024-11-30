CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    description TEXT,
    statusId INT NOT NULL REFERENCES statuses(id) ON DELETE
    SET NULL,
        userId INT REFERENCES users(id) ON DELETE
    SET NULL,
        projectId INT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
        priorityId INT NOT NULL REFERENCES priority(id) ON DELETE CASCADE,
        workspaceId INT References workspaces(id) ON DELETE CASCADE,
        order INT,
        taskProgress INT,
        createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        deletedAt TIMESTAMP,
        deletedBy INT REFERENCES users(id) ON DELETE
    SET NULL
);