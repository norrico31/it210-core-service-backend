package entities

import "time"

type UserStore interface {
	GetUsers() ([]*User, error)
	GetUserById(id int) (*User, error)
	GetUserByEmail(email string) (*User, error)
	CreateUser(UserCreatePayload) error
	UpdateUser(User) error
	DeleteUser(int) error
	SetUserActive(int) error
	UpdateLastActiveTime(int, time.Time) error
}

type User struct {
	ID           int        `json:"id"`
	FirstName    string     `json:"firstName"`
	LastName     string     `json:"lastName"`
	Age          int        `json:"age"`
	Email        string     `json:"e-mail"`
	RoleId       *int       `json:"roleId"`
	Role         Role       `json:"role"`
	Password     string     `json:"-"`
	Projects     []Project  `json:"projects"`
	LastActiveAt *time.Time `json:"lastActiveAt,omitempty"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	DeletedAt    *time.Time `json:"deletedAt,omitempty"`
	DeletedBy    *int       `json:"deletedBy,omitempty"`
}

// TODO REFACTOR
type UserCreatePayload struct {
	FirstName  string    `json:"firstName"`
	LastName   *string   `json:"lastName"`
	Age        *int      `json:"age"`
	Email      string    `validate:"required,email"`
	RoleId     *int      `json:"roleId"`
	ProjectIDS *[]string `json:"projectIds"`

	// MUCH BETTER IF THERE'S DEFAULT PASSWORD
	// Password  string `json:"password" validate:"required,min=3,max=130"`
}

type UserUpdatePayload struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"e-mail" validate:"required,email"`
	Password  string `json:"password,omitempty"` // Optional, for password update
}
