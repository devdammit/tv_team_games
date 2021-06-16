package games

import "time"

type Result struct {
	TeamID    uint8     `json:"id"`
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
}

func NewResult(teamID uint8, name string) *Result {
	return &Result{
		TeamID:    teamID,
		Name:      name,
		Timestamp: time.Now(),
	}
}
