package entities

import (
	"time"
)

type PriorityStore interface {
	GetPriorities() ([]Priority, error)
	GetPriority(int) (*Priority, error)
	CreatePriority(PriorityPayload) (*Priority, error)
	UpdatePriority(PriorityPayload) error
	DeletePriority(int) error
	RestorePriority(int) error
}

type Priority struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
	DeletedBy   *time.Time `json:"deletedBy,omitempty"`
}

type PriorityPayload struct {
	ID          int    `json:"id"`
	Name        string `validate:"required,min=3,max=50"`
	Description string `json:"description"`
}
