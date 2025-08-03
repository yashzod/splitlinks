package model

import (
	"time"

	"github.com/google/uuid"
)

type Experiment struct {
	ID        uuid.UUID         `db:"id"`
	Slug      string            `db:"slug"`
	Metadata  map[string]string `db:"metadata"`
	CreatedAt time.Time         `db:"created_at"`
	Name      string            `db:"name"`
}

type Variant struct {
	ID           uuid.UUID           `db:"id"`
	ExperimentID uuid.UUID           `db:"experiment_id"`
	Name         string              `db:"name"`
	URL          string              `db:"url"`
	Weight       int                 `db:"weight"`
	Targeting    map[string][]string `db:"targeting"`
}
