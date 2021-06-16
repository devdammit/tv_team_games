package games

import "time"

const StartedPoint = "STARTED"
const StoppedPoint = "STOPPED"

type PointAction struct {
	Type string `json:"type"`
	Number int `json:"number"`
	Results []Result `json:"results"`
}


type TeamVoteAction struct {
	ID int `json:"id"`
	Name string `json:"name"`
	RangeTime time.Time `json:"range_time"`
}