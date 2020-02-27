package models

import "time"

// Base fields for a typical DB table
type Base struct {
	ID uint64
	CreatedAt time.Time
	ModifiedAt time.Time
}
