CREATE TABLE IF NOT EXISTS users_projects (
    user_id INT REFERENCES users(id),
    project_id INT REFERENCES projects(id),
    deletedAt TIMESTAMP,
    deletedBy INT REFERENCES users(id),
    PRIMARY KEY (user_id, project_id)
);