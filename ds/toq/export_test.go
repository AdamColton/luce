package toq

import "time"

// SetClock replaces the clock the TimeoutQueue uses so that tests can control
// time. It must be called before anything is added to the queue.
func (tq *TimeoutQueue) SetClock(now func() time.Time, sleep func(time.Duration)) {
	tq.now, tq.sleep = now, sleep
}
