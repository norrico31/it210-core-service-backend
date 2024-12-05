package entities

import "time"

type ProjectStore interface {
	GetProjects(string) ([]*Project, error)
	GetProject(int) (*Project, error)
	ProjectCreate(ProjectCreatePayload) (map[string]interface{}, error)
	ProjectUpdate(int, ProjectUpdatePayload, []int) error
	ProjectDelete(int) (*Project, error)
	ProjectRestore(int) (*Project, error)
}

type Project struct {
	ID           int            `json:"id"`
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	Progress     *float64       `json:"progress"`
	Url          *string        `json:"url"`
	StatusID     int            `json:"statusId"`
	Status       Status         `json:"status"`
	SegmentID    int            `json:"segmentId"`
	Segment      Segment        `json:"segment"`
	DateStarted  *time.Time     `json:"dateStarted"`
	DateDeadline *time.Time     `json:"dateDeadline"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	Users        []User         `json:"users"`
	DeletedBy    *int           `json:"deletedBy,omitempty"`
	DeletedAt    *time.Time     `json:"deletedAt,omitempty"`
	Tasks        []TasksProject `json:"Tasks"`
}

type ProjectCreatePayload struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Progress     *float64 `json:"progress"`
	Url          *string  `json:"url"`
	StatusID     int      `json:"statusId"`
	SegmentID    *int     `json:"segmentId"`
	DateStarted  string   `json:"dateStarted"`
	DateDeadline string   `json:"dateDeadline"`
	UserIDs      *[]int   `json:"userIds"`
}

type ProjectUpdatePayload struct {
	ID           int      `json:"id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Progress     *float64 `json:"progress"`
	Url          *string  `json:"url"`
	StatusID     *int     `json:"statusId"`
	SegmentID    *int     `json:"segmentId"`
	DateStarted  string   `json:"dateStarted"`
	DateDeadline string   `json:"dateDeadline"`
	UserIDs      *[]int   `json:"userIds"`
}
