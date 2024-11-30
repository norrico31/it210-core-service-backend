CREATE TABLE IF NOT EXISTS projects (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    progress DECIMAL(5, 2) DEFAULT 0.00 CHECK (
        progress >= 0
        AND progress <= 100
    ),
    url TEXT,
    statusId: INT REFERENCES statuses(id) ON DELETE
    SET NULL,
        segmentId: INT REFERENCES segment(id) ON DELETE
    SET NULL,
        dateStarted TIMESTAMP,
        dateDeadline TIMESTAMP,
        createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        deletedAt TIMESTAMP,
        deletedBy INT REFERENCES users(id) ON DELETE
    SET NULL
);