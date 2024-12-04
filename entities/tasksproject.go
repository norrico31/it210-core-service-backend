package entities

import (
	"time"
)

type TasksProjectStore interface {
	GetTasksProject(int) ([]*TasksProject, error)
	GetTaskProject(int) (*TasksProject, error)
	TasksProjectCreate(TasksProjectCreatePayload) (*TasksProject, error)
	TasksProjectUpdate(TasksProjectUpdatePayload) error
	TasksProjectDelete(int) error
	TasksProjectRestore(int) (*TasksProject, error)
}

type TasksProject struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	UserID      *int       `json:"userId"`
	User        User       `json:"user"`
	PriorityID  int        `json:"priorityId"`
	ProjectID   int        `json:"projectId"`
	Project     Project    `json:"project"`
	Priority    Priority   `json:"priority"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt"`
	DeletedBy   *int       `json:"deletedBy"`
}

type TasksProjectCreatePayload struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	PriorityID  int    `json:"priorityId"`
	UserID      int    `json:"userId,omitempty"`
	ProjectID   int    `json:"projectId"`
}

type TasksProjectUpdatePayload struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PriorityID  int    `json:"priorityId"`
	UserID      int    `json:"userId,omitempty"`
	ProjectID   int    `json:"projectId"`
}
