-- Create tasks table
CREATE TABLE IF NOT EXISTS project_tasks (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    userId INT REFERENCES users(id) ON DELETE
    SET NULL,
        priorityId INT NOT NULL REFERENCES priorities(id) ON DELETE CASCADE,
        projectId INT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
        createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        deletedAt TIMESTAMP,
        deletedBy INT REFERENCES users(id) ON DELETE
    SET NULL
);