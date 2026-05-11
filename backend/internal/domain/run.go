package domain

import (
	"time"
)

type RunStatus string

const (
	RunStatusPending RunStatus = "pending"
	RunStatusSuccess RunStatus = "success"
	RunStatusFailed  RunStatus = "failed"
)

type Run struct {
	ID            string
	AssistantID   string
	AssistantName string
	CategoryID    string
	CategoryName  string
	UserID        string
	Model         string
	UserPrompt    string
	Output        string
	Status        RunStatus
	Error         string
	CreatedAt     time.Time
}
