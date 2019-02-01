package winddown

import (
	"fmt"
	"time"
)

// State is a possible condition that a winddown timer can be in.
type State string

// Map holds timer durations for various states.
type Map map[State]time.Duration

// A Timer sends a signal on a channel when an amount of time has passed.
// The amount of time it waits before firing is dependent on its current
// state.
type Timer struct {
	C      <-chan time.Time
	timer  *time.Timer
	states Map
}

// NewTimer creates a winddown timer with an initial state and state map.
//
// If asked to adopt a state not in its map it will panic.
func NewTimer(initial State, states Map) Timer {
	duration, ok := states[initial]
	if !ok {
		panic(fmt.Errorf("winddown: timer created with unknown state \"%s\"", initial))
	}
	timer := time.NewTimer(duration)
	return Timer{
		C:      timer.C,
		timer:  timer,
		states: states,
	}
}

// Stop stops the timer and drains its channel if necessary.
//
// It may not be safe to call Stop while another goroutine is listening on
// the timer's channel.
func (t *Timer) Stop() {
	if !t.timer.Stop() {
		select {
		case <-t.timer.C:
		default:
		}
	}
}

// Reset updates the timer's time remaining to the requested state. The
// time remaining will be reset even if the timer was already in the given
// state.
//
// If asked to adopt a state not in its map it will panic.
func (t *Timer) Reset(state State) {
	duration, ok := t.states[state]
	if !ok {
		panic(fmt.Errorf("winddown: timer reset to unknown state \"%s\"", state))
	}
	t.Stop()
	t.timer.Reset(duration)
}
