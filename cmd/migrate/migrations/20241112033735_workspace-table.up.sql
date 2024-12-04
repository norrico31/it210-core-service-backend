-- Step 1: Create the Workspaces Table
CREATE TABLE IF NOT EXISTS workspaces (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    projectId INT NOT NULL REFERENCES projects(id),
    colOrder INT,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deletedAt TIMESTAMP,
    deletedBy INT REFERENCES users(id) ON DELETE
    SET NULL
);
CREATE OR REPLACE FUNCTION set_col_order() RETURNS TRIGGER AS $$ BEGIN -- Assign colOrder if it's not explicitly provided
    IF NEW.colOrder IS NULL THEN
SELECT COALESCE(MAX(colOrder), 0) + 1 INTO NEW.colOrder
FROM workspaces
WHERE projectId = NEW.projectId;
END IF;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER increment_col_order BEFORE
INSERT ON workspaces FOR EACH ROW EXECUTE FUNCTION set_col_order();