package models

import "time"

type Message struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	IncidentID uint      `json:"incident_id" gorm:"index"`
	Role       string    `json:"role"` // "user" or "assistant"
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}
