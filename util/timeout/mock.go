package timeout

import (
	"sync"
	"time"
)

// MockTime is a fake clock for testing code that reads the time, sleeps or
// uses timers. Time only moves when Tick or Advance is called, so a test
// controls it exactly. Sleep blocks until the clock has been advanced far
// enough, and BlockUntil lets a test wait until the Go routines it started are
// asleep before it moves the clock. NewTimer and BlockUntilTimers do the same
// for timers. Use NewMockTime to create one. A MockTime is threadsafe.
type MockTime struct {
	// TickDuration is how far Tick moves the clock for each tick.
	TickDuration time.Duration

	mux      sync.Mutex
	cond     *sync.Cond
	now      time.Time
	sleepers []sleeper
	timers   []*MockTimer
}

type sleeper struct {
	end  time.Time
	wake chan struct{}
}

// MockTimer is a fake time.Timer created by MockTime.NewTimer.
type MockTimer struct {
	// C receives the time when the timer expires.
	C <-chan time.Time

	mt     *MockTime
	c      chan time.Time
	end    time.Time
	active bool
}

// NewMockTime creates a MockTime that starts at the current time and has a
// TickDuration of a millisecond.
func NewMockTime() *MockTime {
	mt := &MockTime{
		TickDuration: time.Millisecond,
		now:          time.Now(),
	}
	mt.cond = sync.NewCond(&mt.mux)
	return mt
}

// Now returns the mock time. It fulfills the signature of time.Now.
func (mt *MockTime) Now() time.Time {
	mt.mux.Lock()
	defer mt.mux.Unlock()
	return mt.now
}

// Tick moves the clock forward by ticks multiplied by TickDuration.
func (mt *MockTime) Tick(ticks int) {
	mt.Advance(time.Duration(ticks) * mt.TickDuration)
}

// Advance moves the clock forward by d, wakes every Sleep that has now
// finished and expires every timer that is now due. When Advance returns,
// those sleeps have been released and those timers have sent on their channels,
// though the Go routines waiting on them may not have run yet. d should not be
// negative.
func (mt *MockTime) Advance(d time.Duration) {
	mt.mux.Lock()
	defer mt.mux.Unlock()
	mt.now = mt.now.Add(d)

	asleep := mt.sleepers[:0]
	for _, s := range mt.sleepers {
		if s.end.After(mt.now) {
			asleep = append(asleep, s)
		} else {
			close(s.wake)
		}
	}
	mt.sleepers = asleep

	pending := mt.timers[:0]
	for _, t := range mt.timers {
		if t.end.After(mt.now) {
			pending = append(pending, t)
		} else {
			t.active = false
			t.c <- mt.now
		}
	}
	mt.timers = pending

	mt.cond.Broadcast()
}

// Sleep blocks until the clock has been advanced by d. If d is not positive
// it returns immediately. It fulfills the signature of time.Sleep.
func (mt *MockTime) Sleep(d time.Duration) {
	mt.mux.Lock()
	end := mt.now.Add(d)
	if !end.After(mt.now) {
		mt.mux.Unlock()
		return
	}
	wake := make(chan struct{})
	mt.sleepers = append(mt.sleepers, sleeper{
		end:  end,
		wake: wake,
	})
	mt.cond.Broadcast()
	mt.mux.Unlock()
	<-wake
}

// Sleeping returns the number of calls to Sleep that are currently blocked.
func (mt *MockTime) Sleeping() int {
	mt.mux.Lock()
	defer mt.mux.Unlock()
	return len(mt.sleepers)
}

// BlockUntil blocks until at least n calls to Sleep are blocked. Use it to wait
// for Go routines to reach their Sleep before moving the clock.
func (mt *MockTime) BlockUntil(n int) {
	mt.mux.Lock()
	defer mt.mux.Unlock()
	for len(mt.sleepers) < n {
		mt.cond.Wait()
	}
}

// NewTimer creates a timer that expires when the clock has been advanced by d.
// Like time.NewTimer, it sends the time on the timer's channel C, which has a
// buffer of one so the send never blocks. If d is not positive the timer
// expires at once.
func (mt *MockTime) NewTimer(d time.Duration) *MockTimer {
	mt.mux.Lock()
	defer mt.mux.Unlock()
	c := make(chan time.Time, 1)
	t := &MockTimer{
		C:   c,
		mt:  mt,
		c:   c,
		end: mt.now.Add(d),
	}
	if d <= 0 {
		c <- mt.now
		return t
	}
	t.active = true
	mt.timers = append(mt.timers, t)
	mt.cond.Broadcast()
	return t
}

// Timers returns the number of timers that have not expired or been stopped.
func (mt *MockTime) Timers() int {
	mt.mux.Lock()
	defer mt.mux.Unlock()
	return len(mt.timers)
}

// BlockUntilTimers blocks until at least n timers have been created and have not
// expired or been stopped. Use it to wait for Go routines to create their
// timers before moving the clock.
func (mt *MockTime) BlockUntilTimers(n int) {
	mt.mux.Lock()
	defer mt.mux.Unlock()
	for len(mt.timers) < n {
		mt.cond.Wait()
	}
}

// Stop prevents the timer from expiring. Like time.Timer.Stop, it returns true
// if the call stopped the timer and false if the timer had already expired or
// been stopped.
func (t *MockTimer) Stop() bool {
	t.mt.mux.Lock()
	defer t.mt.mux.Unlock()
	if !t.active {
		return false
	}
	t.active = false
	for i, ti := range t.mt.timers {
		if ti == t {
			t.mt.timers = append(t.mt.timers[:i], t.mt.timers[i+1:]...)
			break
		}
	}
	return true
}
