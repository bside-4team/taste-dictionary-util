package entities

import "time"

type Restaurant struct {
	ID           string
	Location     string
	ReferenceURL string
	UpdatedAt    time.Time
	Name         string
}
