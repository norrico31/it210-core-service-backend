package entities

import "time"

type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	Age       int    `json:"age"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	// IsActive     bool       `json:"isActive"` // set this when email confirmation
	LastActiveAt *time.Time `json:"lastActiveAt"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	DeletedAt    *time.Time `json:"deletedAt"`
	Projects     []Project  `json:"projects"`
}

type Role struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description int        `json:"description"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt"`
}

type Project struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description int    `json:"description"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	Password    string `json:"-"`
	// IsActive     bool       `json:"isActive"` // set this when email confirmation
	LastActiveAt *time.Time `json:"lastActiveAt"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	DeletedAt    *time.Time `json:"deletedAt"`
	Users        []User     `json:"users"`
}
