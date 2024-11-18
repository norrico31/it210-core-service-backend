CREATE TABLE IF NOT EXISTS projects (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    progress DECIMAL(5, 2) DEFAULT 0.00,
    dateStarted TIMESTAMP,
    dateDeadline TIMESTAMP,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deletedAt TIMESTAMP,
    deletedBy INT REFERENCES users(id) ON DELETE
    SET NULL
);