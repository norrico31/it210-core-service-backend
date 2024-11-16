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

type UserStore interface {
	GetUsers() ([]*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUserById(id int) (*User, error)
	CreateUser(User) error
	UpdateUser(User) error
	DeleteUser(int) error
	SetUserActive(int) error
	UpdateLastActiveTime(int, time.Time) error
}

// USER MUST HAVE A TEAMS
type UserRegisterPayload struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=3,max=130"`
}

type UserUpdatePayload struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password,omitempty"` // Optional, for password update
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
	Users       []User     `json:"users,omitempty"`
	Tasks       []Task     `json:"tasks,omitempty"`
}
