package winddown_test

import (
	"time"

	"github.com/gentlemanautomaton/winddown"
)

// Timer states
const (
	Active = winddown.State("active")
	Idle   = winddown.State("idle")
)

// Timer state map
var states = winddown.Map{
	Active: time.Millisecond * 5,
	Idle:   time.Hour * 24,
}

func Example() {
	// Create a timer in the idle state
	t := winddown.NewTimer(Idle, states)

	// Always stop the timer when finished
	defer t.Stop()

	// Do some work that keeps the timer in an active state
	t.Reset(Active)
	time.Sleep(time.Millisecond * 50)
	t.Reset(Active)

	// Wait for the timer to expire
	<-t.C
}
