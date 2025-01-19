package dbltimer

import "time"

// NewWithTimer creates a DoubleTimer that gets its timers from newTimer instead
// of time.NewTimer, so that tests can control time. newTimer returns the
// channel the timer expires on and a func that stops the timer.
func NewWithTimer(hard, soft time.Duration, callback func(), newTimer func(time.Duration) (<-chan time.Time, func() bool)) *DoubleTimer {
	return newDoubleTimer(hard, soft, callback, newTimer)
}
