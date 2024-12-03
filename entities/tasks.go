package entities

import (
	"time"
)

type TaskStore interface {
	GetTasks() ([]*Task, error)
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
	UserID      *int       `json:"userId"`
	User        User       `json:"user"`
	PriorityID  int        `json:"priorityId"`
	Priority    Priority   `json:"priority"`
	WorkspaceID int        `json:"workspaceId"`
	Workspace   Workspace  `json:"workspace"`
	TaskOrder   int        `json:"taskOrder"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt"`
	DeletedBy   *int       `json:"deletedBy"`
}

type TaskCreatePayload struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	PriorityID  int    `json:"priorityId"`
	WorkspaceID int    `json:"workspaceId"`
	UserID      int    `json:"userId,omitempty"`
}

type TaskUpdatePayload struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	PriorityID  int    `json:"priorityId"`
	WorkspaceID int    `json:"workspaceId"`
	UserID      int    `json:"userId,omitempty"`
}

// TODO DRAGNDROP for TaskORDER
