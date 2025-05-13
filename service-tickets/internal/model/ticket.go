package model

import "time"

// Ticket represents a service ticket in the system
type Ticket struct {
	ID          string    `json:"id"`
	Status      string    `json:"status"`
	User        string    `json:"user"`
	Agent       *string   `json:"agent,omitempty"`
	ProblemID   *int64    `json:"problemId,omitempty"`
	VerticalID  *int64    `json:"verticalId,omitempty"`
	SkillID     *int64    `json:"skillId,omitempty"`
	UserGroupID *int64    `json:"userGroupId,omitempty"`
	Channel     string    `json:"channel"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
