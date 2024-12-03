package entities

import (
	"time"
)

type SegmentsStore interface {
	GetSegments() ([]Segment, error)
	GetSegment(int) (*Segment, error)
	CreateSegment(SegmentPayload) (*Segment, error)
	UpdateSegment(SegmentPayload) error
	DeleteSegment(int) error
	RestoreSegment(int) error
}

type Segment struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
	DeletedBy   *time.Time `json:"deletedBy,omitempty"`
	Projects    []Project  `json:"projects"`
}

type SegmentPayload struct {
	ID          int    `json:"id"`
	Name        string `validate:"required,min=3,max=50"`
	Description string `json:"description"`
	ProjectIDs  *[]int `json:"projectIds"`
}
