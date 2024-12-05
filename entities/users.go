package entities

import "time"

type UserStore interface {
	Login(UserLoginPayload) (string, User, error)
	GetUsers() ([]*User, error)
	GetUserById(id int) (*User, error)
	GetUserByEmail(email string) (*User, error)
	CreateUser(UserCreatePayload) error
	UpdateUser(int, UserUpdatePayload, []int) error
	DeleteUser(int) error
	RestoreUser(int) error
	SetUserActive(int) error
	UpdateLastActiveTime(int, time.Time) error
}

type User struct {
	ID           int        `json:"id"`
	FirstName    string     `json:"firstName"`
	LastName     string     `json:"lastName"`
	Age          int        `json:"age"`
	Email        string     `json:"email"`
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

type UserLoginPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password,omitempty"` // Optional, for password update
}

type UserCreatePayload struct {
	FirstName  string  `json:"firstName"`
	LastName   *string `json:"lastName"`
	Age        *int    `json:"age"`
	Email      string  `validate:"required,email"`
	RoleId     *int    `json:"roleId"`
	ProjectIDS *[]int  `json:"projectIds"`
	Password   string  `json:"-"`
}

type UserUpdatePayload struct {
	FirstName  *string `json:"firstName"`
	LastName   *string `json:"lastName"`
	Age        *int    `json:"age"`
	Email      *string `validate:"required,email"`
	RoleId     *int    `json:"roleId"`
	Password   *string `json:"-"`
	ProjectIDS *[]int  `json:"projectIds"`
}
