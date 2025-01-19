package dbltimer

import (
	"sync/atomic"
	"time"
)

// DoubleTimer runs two timers in parallel. A hard timer that expires after a
// set time and a soft timer that can be reset. An example use would be saving a
// resource to a persistant store.
type DoubleTimer struct {
	hard, soft *time.Timer
	softD      time.Duration
	reset      chan bool
	Callback   func()
	lock       uint32
}

// New creates a DoubleTimer with provided hard and soft timers. When it
// expires it will invoke callback.
func New(hard, soft time.Duration, callback func()) *DoubleTimer {
	dt := &DoubleTimer{
		hard:     time.NewTimer(hard),
		softD:    soft,
		Callback: callback,

		// sending false does a reset
		// sending true does a cancel
		reset: make(chan bool),
	}
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

// Done checks if the timer has expired
func (dt *DoubleTimer) Done() bool {
	return dt.lock != 0
}

func (dt *DoubleTimer) softReset() {
	dt.soft = time.NewTimer(dt.softD)
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
	dt.hard.Stop()
	dt.soft.Stop()
}

// Cancel the DoubleTimer. The returned bool indicates if the cancel was
// successful. If cancel returns true it is guarenteed that the the callback
// will not be invoked.
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
		case <-dt.hard.C:
			dt.callback(hard)
			return
		case <-dt.soft.C:
			dt.callback(soft)
			return
		case cancel := <-dt.reset:
			if cancel {
				return
			}
		}
	}
}
