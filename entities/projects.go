package entities

import "time"

type ProjectStore interface {
	GetProjects() ([]*Project, error)
}

type Project struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt"`
	Users       []User     `json:"users,omitempty"`
	Tasks       []Task     `json:"tasks,omitempty"`
}
