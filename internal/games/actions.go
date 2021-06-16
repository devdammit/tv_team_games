package games

import "time"

const StartedPoint = "STARTED"
const FinishedPoint = "FINISHED"
const UpdateResultPoint = "UPDATE_RESULT"

type PointAction struct {
	Type    string    `json:"type"`
	Number  int       `json:"number"`
	Results []Result  `json:"results"`
	Date    time.Time `json:"date"`
}

type TeamVoteAction struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	RangeTime time.Time `json:"range_time"`
}
