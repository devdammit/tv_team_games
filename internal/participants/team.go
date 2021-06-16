package participants

import (
	"github.com/stianeikeland/go-rpio"
	"time"
)

type TeamObserver interface {
	GetID() int
	OnTeamPress(team *Team)
}

type TeamOpts struct {
	Name string
	Index int
	ButtonPIN uint8
}

type TeamEvent interface {
	Subscribe(observer TeamObserver)
	UnSubscribe(observer TeamObserver)
}

type Team struct {
	Name string
	Index int
	Device rpio.Pin
	isPressed bool
	observers []TeamObserver
}

func NewTeam(opts TeamOpts) *Team {
	button := rpio.Pin(opts.ButtonPIN)

	button.Input()

	team := Team{
		Name: opts.Name,
		Index: opts.Index,
		Device: button,
		isPressed: false,
	}

	go func() {
		for {
			state := button.Read()

			if (state == rpio.High) && !team.isPressed {
				team.isPressed = true
				team.Emit()
			} else if state == rpio.Low {
				team.isPressed = false
			}

			time.Sleep(time.Millisecond * 50)
		}
	}()

	return &team
}


func (p *Team) Emit() {
	for _, observer := range p.observers {
		observer.OnTeamPress(p)
	}
}

func (p *Team) Subscribe(observer TeamObserver) {
	p.observers = append(p.observers, observer)
}

func (p *Team) UnSubscribe(observer TeamObserver) {
	p.removeFromSlice(observer)
}

func (p *Team) removeFromSlice(target TeamObserver) {
	length := len(p.observers)

	for i, observer := range p.observers {
		if target.GetID() == observer.GetID() {
			p.observers[length - 1], p.observers[i] = p.observers[i], p.observers[length - 1]

			p.observers = p.observers[:length - 1]
		}
	}
}

