package entities

import "time"

type User struct {
	ID           int        `json:"id"`
	FirstName    string     `json:"firstName"`
	Age          int        `json:"age"`
	LastName     string     `json:"lastName"`
	Email        string     `json:"email"`
	Password     string     `json:"-"`
	LastActiveAt *time.Time `json:"lastActiveAt"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	DeletedAt    *time.Time `json:"deletedAt"`
	Projects     []Project  `json:"projects"`
}

type Status struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt"`
}

type Task struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	SubTask     []string   `json:"subTask"`
	Description string     `json:"description"`
	StatusID    int        `json:"statusId"`
	Status      Status     `json:"status"`
	User        User       `json:"user"`
	UserId      int        `json:"userId"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt"`
	Projects    []Project  `json:"projects"`
}

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
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt"`
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
