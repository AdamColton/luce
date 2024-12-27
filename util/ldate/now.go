package ldate

import (
	"sync"
	"time"
)

var (
	nowMu sync.RWMutex
	now   Date
)

func init() {
	setNow(TimeToDate(time.Now()))
	go updateNow(time.Now, nil)
}

func setNow(d Date) {
	nowMu.Lock()
	now = d
	nowMu.Unlock()
}

// Now returns today's Date in the local time zone. It is updated just after
// midnight.
func Now() Date {
	nowMu.RLock()
	defer nowMu.RUnlock()
	return now
}

// updateNow keeps now current: it sets it from clock, waits for the next day and
// repeats until stop is closed. A nil stop is never closed.
func updateNow(clock func() time.Time, stop <-chan struct{}) {
	for {
		n := clock()
		setNow(TimeToDate(n))
		tomorrow := time.Date(n.Year(), n.Month(), n.Day()+1, 0, 0, 0, 100, n.Location())
		timer := time.NewTimer(tomorrow.Sub(n))
		select {
		case <-timer.C:
		case <-stop:
			timer.Stop()
			return
		}
	}
}
