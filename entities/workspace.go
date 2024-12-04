package entities

import (
	"time"
)

type WorkspaceStore interface {
	GetWorkspaces() ([]Workspace, error)
	GetWorkspace(int) ([]Workspace, error)
	CreateWorkspace(WorkspacePayload) (*Workspace, error)
	UpdateWorkspace(WorkspacePayload) error
	DeleteWorkspace(int) error
	RestoreWorkspace(int) error
	TaskDragNDrop(int, int, int) error
}

type Workspace struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	ProjectID   int        `json:"projectId"`
	Project     Project    `json:"project"`
	ColOrder    int        `json:"colOrder"`
	CreatedAt   time.Time  `json:"createdAt"`
	Tasks       []Task     `json:"tasks"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
	DeletedBy   *time.Time `json:"deletedBy,omitempty"`
}

type WorkspacePayload struct {
	ID          int    `json:"id"`
	Name        string `validate:"required,min=3,max=50"`
	Description string `json:"description"`
	ProjectID   int    `json:"projectId"`
	ColOrder    int    // Optional
}

type TaskDragNDrop struct {
	SourceTaskId      int `json:"sourceTaskId"`
	DestinationTaskId int `json:"destinationTaskId"`
}

// TODO DRAGNDROP for colOrder
