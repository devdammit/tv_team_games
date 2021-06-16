package main

import (
	"fmt"
	"github.com/r3labs/sse/v2"
	"github.com/stianeikeland/go-rpio"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"tv_vote_team/internal/games"
	"tv_vote_team/internal/participants"
)

func WaitForCtrlC() {
	var endWaiter sync.WaitGroup
	endWaiter.Add(1)
	var signalChannel chan os.Signal
	signalChannel = make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	go func() {
		<-signalChannel
		endWaiter.Done()
	}()
	endWaiter.Wait()
}

//
//


func main() {
	err := rpio.Open()
	if err != nil {
		panic(fmt.Sprint("unable to open gpio", err.Error()))
	}

	defer rpio.Close()

	gameOpts := games.Options{
		RefereeButtonPIN: 17,
		Teams: []participants.TeamOpts{
			{
				ButtonPIN: 27,
				Name: "Cepheus",
			},
			{
				ButtonPIN: 23,
				Name: "Draco",
			},
			{
				ButtonPIN: 24,
				Name: "Hercules",
			},
		},
	}

	game := games.NewGame(gameOpts)

	fmt.Println("Starting game")
	game.Start()


	server := sse.New()
	server.CreateStream("point")
	server.CreateStream("vote")
	mux := http.NewServeMux()
	mux.HandleFunc("/events", server.HTTPHandler)

	game.Server = server

	http.ListenAndServe(":8080", mux)


	WaitForCtrlC()
}
