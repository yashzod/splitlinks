package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type User struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	Email      string         `gorm:"uniqueIndex;not null" json:"email"`
	FullName   string         `gorm:"size:255" json:"full_name"`
	FirstName  string         `gorm:"size:100" json:"first_name"`
	LastName   string         `gorm:"size:100" json:"last_name"`
	ImageURL   string         `gorm:"size:1024" json:"image_url"`
	IsDisabled bool           `gorm:"default:false" json:"is_disabled"`
	Metadata   datatypes.JSON `gorm:"type:jsonb" json:"metadata,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}
