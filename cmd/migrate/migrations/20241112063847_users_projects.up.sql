CREATE TABLE IF NOT EXISTS users_projects (
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    project_id INT REFERENCES projects(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, project_id)
);