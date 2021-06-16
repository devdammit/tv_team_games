package participants

import (
	"github.com/stianeikeland/go-rpio"
	"time"
)

type RefereeObserver interface {
	GetID() int
	OnRefereePress()
}

type RefereeEvent interface {
	Subscribe(observer RefereeObserver)
	UnSubscribe(observer RefereeObserver)
}



type Referee struct {
	observers []RefereeObserver
	isPressed bool
	Device rpio.Pin
}

func NewReferee(pin uint8) *Referee {
	button := rpio.Pin(pin)

	button.Input()

	referee := Referee{
		Device: button,
		isPressed: false,
	}

	go func() {
		for {
			state := button.Read()

			if (state == rpio.High) && !referee.isPressed {
				referee.isPressed = true
				referee.Emit()
			} else if state == rpio.Low {
				referee.isPressed = false
			}

			time.Sleep(time.Millisecond * 50)
		}
	}()

	return &referee
}



func (p *Referee) Emit() {
	for _, observer := range p.observers {
		observer.OnRefereePress()
	}
}

func (p *Referee) Subscribe(observer RefereeObserver) {
	p.observers = append(p.observers, observer)
}

func (p *Referee) UnSubscribe(observer RefereeObserver) {
	p.removeFromSlice(observer)
}

func (p *Referee) removeFromSlice(target RefereeObserver) {
	length := len(p.observers)

	for i, observer := range p.observers {
		if target.GetID() == observer.GetID() {
			p.observers[length - 1], p.observers[i] = p.observers[i], p.observers[length - 1]

			p.observers = p.observers[:length - 1]
		}
	}
}




