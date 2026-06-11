package models

import (
	"time"

	"gorm.io/gorm"
)

type Severity string
type Status string

const (
	SeverityP1 Severity = "P1"
	SeverityP2 Severity = "P2"
	SeverityP3 Severity = "P3"
	SeverityP4 Severity = "P4"

	StatusOpen       Status = "open"
	StatusInProgress Status = "in_progress"
	StatusResolved   Status = "resolved"
	StatusClosed     Status = "closed"
)

type Incident struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `json:"title" binding:"required"`
	Description string         `json:"description" binding:"required"`
	Severity    Severity       `json:"severity"`
	Status      Status         `json:"status"`
	Analysis    string         `json:"analysis"`
	AnalyzedAt  *time.Time     `json:"analyzed_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type CreateIncidentRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type UpdateIncidentRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Severity    Severity `json:"severity"`
	Status      Status   `json:"status"`
}
