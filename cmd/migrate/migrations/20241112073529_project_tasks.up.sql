CREATE TABLE IF NOT EXISTS project_tasks (
    project_id INT REFERENCES projects(id) ON DELETE CASCADE,
    task_id INT REFERENCES tasks(id) ON DELETE CASCADE,
    PRIMARY KEY (project_id, task_id)
);