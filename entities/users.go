package entities

import "time"

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
	DeletedBy    *int       `json:"deletedBy"`
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
