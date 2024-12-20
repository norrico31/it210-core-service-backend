package entities

import (
	"time"
)

type RoleStore interface {
	GetRoles() ([]Role, error)
	GetRole(int) (*Role, error)
	CreateRole(RolePayload) (*Role, error)
	UpdateRole(RolePayload) error
	DeleteRole(int) error
	RestoreRole(int) error
}

type Role struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
	DeletedBy   *time.Time `json:"deletedBy,omitempty"`
}

type RolePayload struct {
	ID          int    `json:"id"`
	Name        string `validate:"required,min=3,max=50"`
	Description string `json:"description"`
}
