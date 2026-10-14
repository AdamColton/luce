package timer

import (
	"sync"
	"time"

	"github.com/adamcolton/luce/lerr"
)

// State is the stage of a Timer's life.
type State byte

const (
	// Undefined is the State of a Timer that was not created by New or Delay.
	Undefined State = iota
	// Running is the State of a Timer that is counting down to its end.
	Running
	// Paused is the State of a Timer that Pause stopped until Resume is
	// called.
	Paused
	// Done is the State of a Timer that has expired and invoked its callback.
	Done
)

// Timer invokes a callback at a set time. Unlike a time.Timer, it can be
// paused and resumed, and its end can be moved while it runs. A Timer is
// threadsafe.
type Timer struct {
	mux       sync.Mutex
	callback  func()
	now       func() time.Time
	newTimer  newTimerFunc
	end       time.Time
	remaining time.Duration
	state     State
	// stop is closed to end the Go routine waiting on the current run.
	stop      chan struct{}
	stopTimer func() bool
}

// ErrNilCallback is the panic value when New or Delay is given a nil callback.
const ErrNilCallback = lerr.Str("callback cannot be nil")

// Delay creates a Running Timer that invokes callback once d has passed. It
// panics with ErrNilCallback if callback is nil.
func Delay(d time.Duration, callback func()) *Timer {
	return New(time.Now().Add(d), callback)
}

// New creates a Running Timer that invokes callback at end. It panics with
// ErrNilCallback if callback is nil.
func New(end time.Time, callback func()) *Timer {
	return newTimerWithClock(end, callback, time.Now, newRealTimer)
}

func newTimerWithClock(end time.Time, callback func(), now func() time.Time, newTimer newTimerFunc) *Timer {
	if callback == nil {
		panic(ErrNilCallback)
	}
	t := &Timer{
		callback: callback,
		now:      now,
		newTimer: newTimer,
	}
	t.Update(end)
	return t
}

// start runs the Timer for d. The caller must hold mux.
func (t *Timer) start(d time.Duration) {
	c, stopTimer := t.newTimer(d)
	stop := make(chan struct{})
	t.stop, t.stopTimer = stop, stopTimer
	t.state = Running
	go func() {
		select {
		case <-stop:
		case <-c:
			t.expire(stop)
		}
	}()
}

// expire ends the run that stop belongs to and invokes the callback, unless
// Update or Pause stopped that run while the expiring Go routine waited for
// mux.
func (t *Timer) expire(stop chan struct{}) {
	t.mux.Lock()
	select {
	case <-stop:
		t.mux.Unlock()
		return
	default:
	}
	t.state = Done
	t.mux.Unlock()
	t.callback()
}

// halt stops the current run. The caller must hold mux and the Timer must be
// Running.
func (t *Timer) halt() {
	close(t.stop)
	t.stopTimer()
}

// Update sets the Timer to expire at end and leaves it Running, whether it
// was Running or Paused. It does nothing if the Timer is Done.
func (t *Timer) Update(end time.Time) {
	t.mux.Lock()
	defer t.mux.Unlock()
	switch t.state {
	case Done:
		return
	case Running:
		t.halt()
	}
	t.end = end
	t.remaining = 0
	t.start(end.Sub(t.now()))
}

// UpdateDelay sets the Timer to expire once d has passed. A Running Timer
// expires d from now. A Paused Timer stays Paused and expires d after it is
// resumed. It does nothing if the Timer is Done.
func (t *Timer) UpdateDelay(d time.Duration) {
	t.mux.Lock()
	defer t.mux.Unlock()
	switch t.state {
	case Running:
		t.halt()
		t.end = t.now().Add(d)
		t.start(d)
	case Paused:
		t.remaining = d
	}
}

// Pause stops a Running Timer, keeping the time it has left until Resume is
// called. It does nothing if the Timer is not Running.
func (t *Timer) Pause() {
	t.mux.Lock()
	defer t.mux.Unlock()
	if t.state != Running {
		return
	}
	t.halt()
	t.state = Paused
	t.remaining = t.end.Sub(t.now())
}

// Resume restarts a Paused Timer with the time it had left when it was
// paused. It does nothing if the Timer is not Paused.
func (t *Timer) Resume() {
	t.mux.Lock()
	defer t.mux.Unlock()
	if t.state != Paused {
		return
	}
	t.end = t.now().Add(t.remaining)
	t.start(t.remaining)
	t.remaining = 0
}

// State returns the Timer's current State.
func (t *Timer) State() State {
	t.mux.Lock()
	defer t.mux.Unlock()
	return t.state
}

// End returns when the Timer expires. For a Paused Timer, that is when it would
// expire if it were resumed now.
func (t *Timer) End() time.Time {
	t.mux.Lock()
	defer t.mux.Unlock()
	if t.state == Paused {
		return t.now().Add(t.remaining)
	}
	return t.end
}
