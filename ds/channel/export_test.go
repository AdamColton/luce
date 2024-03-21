package channel

import "time"

// SetTimer replaces the timer that ends each Cycle, so that a test decides when
// a Cycle ends. newTimer is called at the start of every Cycle with the two
// delays and the callback to invoke when the timer expires. What it returns is
// Reset each time more data arrives.
func (m *Merge[T]) SetTimer(newTimer func(hard, soft time.Duration, callback func()) interface{ Reset() bool }) {
	m.newTimer = func(hard, soft time.Duration, callback func()) timer {
		return newTimer(hard, soft, callback)
	}
}

// SetTimeout replaces the Timeout that Cycle uses to wait for the timer once
// the incoming channel is closed.
func (m *Merge[T]) SetTimeout(timeout func(d time.Duration, ch <-chan bool) (bool, error)) {
	m.timeout = timeout
}
