package models

import (
	"time"
)

type (
	CreationTime struct {
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	TableConfig struct {
		Name string
	}
)

// UpdateTimes - Update the time variables that need to be addressed when making changes
func (m *CreationTime) UpdateTimes() {
	now := time.Now()

	if m.CreatedAt.IsZero() {
		m.CreatedAt = now
	}

	m.UpdatedAt = now
}
