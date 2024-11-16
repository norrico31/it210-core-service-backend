package entities

import "time"

type ProjectStore interface {
	GetProjects() ([]*Project, error)
	GetProject(int) (*Project, error)
	ProjectCreate(ProjectCreatePayload) (*Project, error)
	ProjectUpdate(ProjectUpdatePayload) (*Project, error)
}

type Project struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt"`
	Users       []User     `json:"users"`
	Tasks       []Task     `json:"tasks"`
}

type ProjectCreatePayload struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ProjectUpdatePayload struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
