CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    description TEXT,
    SET NULL,
    userId INT REFERENCES users(id) ON DELETE
    SET NULL,
        projectId INT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
        priorityId INT NOT NULL REFERENCES priorities(id) ON DELETE CASCADE,
        workspaceId INT REFERENCES workspaces(id) ON DELETE CASCADE,
        taskOrder INT,
        createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        deletedAt TIMESTAMP,
        deletedBy INT REFERENCES users(id) ON DELETE
    SET NULL
);
CREATE OR REPLACE FUNCTION set_task_order() RETURNS TRIGGER AS $$ BEGIN IF NEW.taskOrder IS NULL THEN
SELECT COALESCE(MAX(taskOrder), 0) + 1 INTO NEW.taskOrder
FROM tasks
WHERE workspaceId = NEW.workspaceId;
END IF;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER increment_task_order BEFORE
INSERT ON tasks FOR EACH ROW EXECUTE FUNCTION set_task_order();