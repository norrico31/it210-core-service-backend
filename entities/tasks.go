package entities

import "time"

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
