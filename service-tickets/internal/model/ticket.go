package model

import "time"

// Ticket represents a service ticket in the system
type Ticket struct {
	ID         string    `json:"id"`
	VerticalID string    `json:"verticalId"`
	UserID     string    `json:"userId"`
	Assign     string    `json:"assign"`
	SkillID    string    `json:"skillId"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}
