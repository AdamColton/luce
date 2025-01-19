// Package dbltimer provides a timer that invokes a callback when it expires.
// The timer can be extended, but never past a hard limit.
package dbltimer

import (
	"sync/atomic"
	"time"
)

// DoubleTimer runs two timers in parallel: a hard timer that expires after a
// set time and a soft timer that can be reset. The DoubleTimer expires when
// either timer does, and then Callback is invoked. Callback is invoked at most
// once, and never if Cancel succeeded. An example use would be saving a
// resource to a persistent store.
//
// Reset and Cancel can block forever if they are called while the DoubleTimer
// is expiring or being canceled from another Go routine.
type DoubleTimer struct {
	hard, soft timer
	newTimer   newTimerFunc
	softD      time.Duration
	reset      chan bool
	Callback   func()
	lock       uint32
}

// timer is the part of a time.Timer that a DoubleTimer uses.
type timer struct {
	c    <-chan time.Time
	stop func() bool
}

// newTimerFunc creates a timer that expires after d. It returns the channel the
// timer expires on and a func that stops the timer.
type newTimerFunc func(d time.Duration) (c <-chan time.Time, stop func() bool)

// newRealTimer creates a time.Timer. It is what a DoubleTimer uses except in
// tests, which replace it to control time.
func newRealTimer(d time.Duration) (<-chan time.Time, func() bool) {
	t := time.NewTimer(d)
	return t.C, t.Stop
}

// New creates a DoubleTimer with provided hard and soft timers. When it
// expires it will invoke callback.
func New(hard, soft time.Duration, callback func()) *DoubleTimer {
	return newDoubleTimer(hard, soft, callback, newRealTimer)
}

func newDoubleTimer(hard, soft time.Duration, callback func(), newTimer newTimerFunc) *DoubleTimer {
	dt := &DoubleTimer{
		newTimer: newTimer,
		softD:    soft,
		Callback: callback,

		// sending false does a reset
		// sending true does a cancel
		reset: make(chan bool),
	}
	dt.hard.c, dt.hard.stop = newTimer(hard)
	go dt.run()
	return dt
}

func (dt *DoubleTimer) drainReset() chan<- bool {
	complete := make(chan bool)
	go func() {
		for {
			select {
			case <-dt.reset:
				// do nothing, just drain the channel
			case <-complete:
				return

			}
		}
	}()
	return complete
}

const (
	zero uint32 = iota
	cancel
	soft
	hard
)

func (dt *DoubleTimer) callback(src uint32) {
	if !atomic.CompareAndSwapUint32(&(dt.lock), 0, src) {
		return
	}
	complete := dt.drainReset()
	dt.Callback()
	complete <- true
}

// Done reports if the timer has expired or been canceled.
func (dt *DoubleTimer) Done() bool {
	return dt.lock != 0
}

func (dt *DoubleTimer) softReset() {
	dt.soft.c, dt.soft.stop = dt.newTimer(dt.softD)
}

// Reset the soft timer. The returned bool indicates if the reset was
// successful. Note that the bool only indicates that the DoubleTimer had
// not expired when the method was invoked, it could expire while the
// method is executing.
func (dt *DoubleTimer) Reset() bool {
	if dt.lock != 0 {
		return false
	}
	dt.softReset()
	dt.reset <- false

	return true
}

func (dt *DoubleTimer) clearTimers() {
	dt.hard.stop()
	dt.soft.stop()
}

// Cancel the DoubleTimer. The returned bool indicates if the cancel was
// successful. If cancel returns true it is guaranteed that the callback will
// not be invoked. It returns false if the DoubleTimer has already expired or
// been canceled.
func (dt *DoubleTimer) Cancel() bool {
	didCancel := atomic.CompareAndSwapUint32(&(dt.lock), 0, cancel)
	if didCancel {
		dt.reset <- true
		dt.clearTimers()
	}

	return didCancel
}

func (dt *DoubleTimer) run() {
	dt.softReset()
	for {
		select {
		case <-dt.hard.c:
			dt.callback(hard)
			return
		case <-dt.soft.c:
			dt.callback(soft)
			return
		case cancel := <-dt.reset:
			if cancel {
				return
			}
		}
	}
}
