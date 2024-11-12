CREATE TABLE IF NOT EXISTS statuses (
    id SERIAL PRIMARY KEY,
    -- Auto-incrementing ID
    name VARCHAR(255) NOT NULL,
    -- Name of the status
    description TEXT,
    -- Description of the status
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    -- Creation timestamp
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    -- Last updated timestamp
    deletedAt TIMESTAMP -- Soft delete timestamp (nullable)
);