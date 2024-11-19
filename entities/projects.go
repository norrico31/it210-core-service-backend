package entities

import "time"

type ProjectStore interface {
	GetProjects(string) ([]*Project, error)
	GetProject(int) (*Project, error)
	ProjectCreate(ProjectCreatePayload) (*Project, error)
	ProjectUpdate(ProjectUpdatePayload) (*Project, error)
	ProjectDelete(int) (*Project, error)
	ProjectRestore(int) (*Project, error)
}

type Project struct {
	ID           int        `json:"id"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	Progress     *float64   `json:"progress"`
	DateStarted  *time.Time `json:"dateStarted"`
	DateDeadline *time.Time `json:"dateDeadline"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	DeletedAt    *time.Time `json:"deletedAt"`
	Users        []User     `json:"users"`
	Tasks        []Task     `json:"tasks"`
	DeletedBy    *int       `json:"deletedBy"`
}

type ProjectCreatePayload struct {
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	Progress     *float64   `json:"progress"`
	DateStarted  *time.Time `json:"dateStarted"`
	DateDeadline *time.Time `json:"dateDeadline"`
}

type ProjectUpdatePayload struct {
	ID           int        `json:"id"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	Progress     *float64   `json:"progress"`
	DateStarted  *time.Time `json:"dateStarted"`
	DateDeadline *time.Time `json:"dateDeadline"`
}
