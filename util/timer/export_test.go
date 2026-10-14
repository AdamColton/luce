package timer

import "time"

// NewDoubleWithTimer creates a DoubleTimer that gets its timers from newTimer
// instead of time.NewTimer, so that tests can control time. newTimer returns
// the channel the timer expires on and a func that stops the timer.
func NewDoubleWithTimer(hard, soft time.Duration, callback func(), newTimer func(time.Duration) (<-chan time.Time, func() bool)) *DoubleTimer {
	return newDoubleTimer(hard, soft, callback, newTimer)
}

// NewWithClock creates a Timer that reads the time from now and gets its timers
// from newTimer instead of the time package, so that tests can control time.
// newTimer returns the channel the timer expires on and a func that stops the
// timer.
func NewWithClock(end time.Time, callback func(), now func() time.Time, newTimer func(time.Duration) (<-chan time.Time, func() bool)) *Timer {
	return newTimerWithClock(end, callback, now, newTimer)
}
