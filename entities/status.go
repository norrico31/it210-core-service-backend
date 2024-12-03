package entities

import "time"

type StatusStore interface {
	GetStatuses() ([]Status, error)
	GetStatus(int) (*Status, error)
	CreateStatus(StatusPayload) (*Status, error)
	UpdateStatus(StatusPayload) error
	DeleteStatus(int) error
	RestoreStatus(int) error
}

type Status struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
	DeletedBy   *time.Time `json:"deletedBy,omitempty"`
}

type StatusPayload struct {
	ID          int    `json:"id"`
	Name        string `validate:"required,min=3,max=50"`
	Description string `json:"description"`
}
