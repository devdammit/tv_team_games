package games

import (
	"fmt"
	"time"
	"tv_vote_team/internal/participants"
)

type EventCallback func()

type ModeType int

const (
	Stopped ModeType = iota
	Started
	Finished
)

type PointObserver interface {
	GetID() string
	OnStartPoint(point *Point)
	OnVoteTeam(team *participants.Team)
	OnFinishPoint(point *Point)
}

type PointEmmiter interface {
	Subscribe(observer PointObserver)
	UnSubscribe(observer PointObserver)
}

type Point struct {
	Number   int
	OnStart  func(cb EventCallback)
	OnFinish func(cb EventCallback)
	Results  map[uint8]Result

	Started   time.Time
	observers []PointObserver
	teamCount int
	mode      ModeType
}

func NewPoint(game *Game, num int) *Point {
	point := Point{
		Number:    num,
		Results:   map[uint8]Result{},
		Started:   time.Now(),
		teamCount: len(game.Teams),
	}

	game.Referee.Subscribe(&point)

	for _, team := range game.Teams {
		team.Subscribe(&point)
	}

	return &point
}

func (p *Point) GetID() int {
	return p.Number
}

func (p *Point) EmitStart() {
	for _, observer := range p.observers {
		observer.OnStartPoint(p)
	}
}

func (p *Point) EmitFinish() {
	for _, observer := range p.observers {
		observer.OnFinishPoint(p)
	}
}

func (p *Point) EmitVoteTeam(team *participants.Team) {
	for _, observer := range p.observers {
		observer.OnVoteTeam(team)
	}
}

func (p *Point) Subscribe(observer PointObserver) {
	p.observers = append(p.observers, observer)
}

func (p *Point) UnSubscribe(observer PointObserver) {
	p.removeFromSlice(observer)
}

func (p *Point) OnRefereePress() {
	switch p.mode {
	case Stopped:
		p.mode = Started
		p.EmitStart()

	case Started:
		p.mode = Finished
		p.EmitFinish()
	}
}

func (p *Point) OnTeamPress(team *participants.Team) {
	if p.mode == Started {
		if _, ok := p.Results[uint8(team.Index)]; !ok {
			result := NewResult(uint8(team.Index))
			p.Results[uint8(team.Index)] = *result
			fmt.Printf("Team #%v pressed \n", team.Index)

			p.EmitVoteTeam(team)

			if len(p.Results) == p.teamCount {
				p.EmitFinish()
			}
		}
	}
}

func (p *Point) removeFromSlice(target PointObserver) {
	length := len(p.observers)

	for i, observer := range p.observers {
		if target.GetID() == observer.GetID() {
			p.observers[length-1], p.observers[i] = p.observers[i], p.observers[length-1]

			p.observers = p.observers[:length-1]
		}
	}
}
