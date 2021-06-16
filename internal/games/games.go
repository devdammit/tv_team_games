package games

import (
	"encoding/json"
	"fmt"
	"github.com/r3labs/sse/v2"
	"time"
	"tv_vote_team/internal/participants"
)

type Game struct {
	Referee *participants.Referee
	Teams   []*participants.Team
	Points  []*Point
	Server  *sse.Server

	CurrentPoint int
}

type Options struct {
	RefereeButtonPIN uint8
	Teams            []participants.TeamOpts
}

func NewGame(opts Options) *Game {
	referee := participants.NewReferee(opts.RefereeButtonPIN)
	teams := make([]*participants.Team, len(opts.Teams))

	for i, teamOpts := range opts.Teams {
		teamOpts.Index = i

		team := participants.NewTeam(teamOpts)

		teams[i] = team
	}

	return &Game{
		Referee: referee,
		Teams:   teams,
	}
}

func (game *Game) Start() {
	game.NextPoint()
}

func (game *Game) NextPoint() {
	point := NewPoint(game, game.CurrentPoint)

	game.Points = append(game.Points, point)
	point.Subscribe(game)
}

func (game *Game) GetID() string {
	return "game"
}

func (game *Game) OnStartPoint(point *Point) {
	fmt.Printf("Started #%v point \n", game.CurrentPoint)

	b, _ := json.Marshal(PointAction{
		Type:   StartedPoint,
		Number: point.Number,
		Date:   time.Now(),
	})

	game.Server.Publish("point", &sse.Event{
		Data: b,
	})
}

func (game *Game) OnFinishPoint(point *Point) {
	fmt.Printf("Finished #%v point \n", game.CurrentPoint)

	var results []Result

	for _, result := range point.Results {
		results = append(results, result)
	}

	b, _ := json.Marshal(PointAction{
		Type:    FinishedPoint,
		Number:  point.Number,
		Results: results,
		Date:    time.Now(),
	})

	game.Server.Publish("point", &sse.Event{
		Data: b,
	})

	game.Points[game.CurrentPoint].UnSubscribe(game)
	game.CurrentPoint++

	game.NextPoint()
}

func (game *Game) OnVoteTeam(team *participants.Team) {
	b, _ := json.Marshal(TeamVoteAction{
		ID:        team.Index,
		Name:      team.Name,
		RangeTime: time.Now(),
	})

	var results []Result

	for _, result := range game.Points[game.CurrentPoint].Results {
		results = append(results, result)
	}

	p, _ := json.Marshal(PointAction{
		Type:    UpdateResultPoint,
		Number:  game.CurrentPoint,
		Results: results,
		Date:    time.Now(),
	})

	game.Server.Publish("vote", &sse.Event{
		Data: b,
	})

	game.Server.Publish("point", &sse.Event{
		Data: p,
	})
}
