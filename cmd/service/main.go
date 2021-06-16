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
		RefereeButtonPIN: 21,
		Teams: []participants.TeamOpts{
			{
				ButtonPIN: 20, // PIN 28
				Name:      "Cepheus",
			},
			{
				ButtonPIN: 26, // pin 25
				Name:      "Draco",
			},
			{
				ButtonPIN: 16, // pin 27
				Name:      "Hercules",
			},
			{
				ButtonPIN: 12, // pin 26
				Name:      "Hydra",
			},
			{
				ButtonPIN: 18, // pin 1
				Name:      "Leo",
			},
			{
				ButtonPIN: 25, // 6
				Name:      "Lyra",
			},
			{
				ButtonPIN: 24, // 5
				Name:      "Perseus",
			},
			{
				ButtonPIN: 22, // 3
				Name:      "Phoenix",
			},
			{
				ButtonPIN: 27, // 2
				Name:      "Saggita",
			},
			{
				ButtonPIN: 17, // 0
				Name:      "Vela",
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
	mux.HandleFunc("/restart", func(w http.ResponseWriter, r *http.Request) {
		game.Start()
		server.Close()
		server.CreateStream("point")
		server.CreateStream("vote")
	})

	game.Server = server

	http.ListenAndServe(":8080", mux)

	WaitForCtrlC()
}
