CREATE TABLE IF NOT EXISTS segments_projects (
    segmentId INT REFERENCES segments(id),
    projectId INT REFERENCES projects(id),
    deletedAt TIMESTAMP,
    deletedBy INT REFERENCES users(id),
    PRIMARY KEY (segmentId, projectId)
);