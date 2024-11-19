package entities

import (
	"time"
)

type RoleStore interface {
	GetRoles() ([]*Role, error)
	GetRole(int) (*Role, error)
	CreateRole(Role) error
	UpdateRole(Role) error
	DeleteRole(int) error
	RestoreRole(int) error
}

type Role struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"createdAt,omitempty"`
	UpdatedAt   time.Time  `json:"updatedAt,omitempty"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
}
