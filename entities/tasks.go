package entities

import (
	"time"
)

type TaskStore interface {
	GetTasks(string) ([]*Task, error)
	GetTask(int) (*Task, error)
	TaskCreate(TaskCreatePayload) (*Task, error)
	TaskUpdate(TaskUpdatePayload) (*Task, error)
	TaskDelete(int) (*Task, error)
	TaskRestore(int) (*Task, error)
}

type Task struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	StatusID    int        `json:"statusId"`
	Status      Status     `json:"status"`
	UserID      *int       `json:"userId"`
	User        User       `json:"user"`
	ProjectID   int        `json:"projectId"`
	Project     Project    `json:"project"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt"`
	DeletedBy   *int       `json:"deletedBy"`
}

type TaskCreatePayload struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	StatusID    int    `json:"statusId,omitempty"`
	UserID      int    `json:"userId,omitempty"`
	ProjectID   int    `json:"projectId,omitempty"`
}

type TaskUpdatePayload struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	StatusID    int    `json:"statusId,omitempty"`
	UserID      int    `json:"userId,omitempty"`
	ProjectID   int    `json:"projectId,omitempty"`
}
