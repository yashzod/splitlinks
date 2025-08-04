package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Experiment struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Slug      string         `gorm:"uniqueIndex"`
	Metadata  datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt time.Time
	Name      string
}

type Variant struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	ExperimentID uuid.UUID `gorm:"type:uuid;index"`
	Name         string
	URL          string
	Weight       int
	Targeting    datatypes.JSON `gorm:"type:jsonb"`
}
