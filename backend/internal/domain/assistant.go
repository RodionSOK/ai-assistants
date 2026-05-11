package domain

import (
	"time"
)

type Assistant struct {
	ID                string
	CategoryID        string
	CategoryName      string
	Name              string
	Description       string
	Model             string
	SystemPrompt      string
	ExampleUserPrompt string
	IsActive          bool
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
