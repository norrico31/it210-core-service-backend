package entities

import (
	"time"

	"github.com/lib/pq"
)

type TaskStore interface {
	GetTasks() ([]*Task, error)
	GetTask(int) (*Task, error)
	TaskCreate(TaskCreatePayload) (*Task, error)
}

// TODO: IN DB TABLE MAKE THE USERID AND STATUSID NULLABLE AND PROJECTID
type Task struct {
	ID          int            `json:"id"`
	Title       string         `json:"title"`
	SubTask     pq.StringArray `json:"subTask"`
	Description string         `json:"description"`
	StatusID    int            `json:"statusId"`
	Status      Status         `json:"status"`
	UserID      *int           `json:"userId"`
	User        User           `json:"user"`
	ProjectID   int            `json:"projectId"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   *time.Time     `json:"deletedAt"`
}

type TaskCreatePayload struct {
	Title       string         `json:"title"`
	SubTask     pq.StringArray `json:"subTask"`
	Description string         `json:"description"`
	StatusID    int            `json:"statusId,omitempty"`
	UserID      int            `json:"userId,omitempty"`
	ProjectID   int            `json:"projectId,omitempty"`
}
