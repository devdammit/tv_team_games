package games

import "time"

type Result struct {
	TeamID uint8 `json:"team_id"`
	Timestamp time.Time `json:"timestamp"`
}

func NewResult(teamID uint8) *Result {
	return &Result{
		TeamID: teamID,
		Timestamp: time.Now(),
	}
}

