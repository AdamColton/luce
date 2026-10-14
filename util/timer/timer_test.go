package timer_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/adamcolton/luce/util/timeout"
	"github.com/adamcolton/luce/util/timer"
	"github.com/stretchr/testify/assert"
)

// Most of these tests use timeout.MockTime as the clock, as the DoubleTimer
// tests do, so they never wait on real time.

// newTimer creates a Timer that uses a MockTime as its clock and expires after
// d. The callback sends on called, which has a buffer so that the callback
// never blocks.
func newTimer(t *testing.T, d time.Duration) (*timer.Timer, *timeout.MockTime, chan bool) {
	mt := timeout.NewMockTime()
	called := make(chan bool, 1)
	tmr := timer.NewWithClock(mt.Now().Add(d), func() {
		called <- true
	}, mt.Now, func(d time.Duration) (<-chan time.Time, func() bool) {
		mtmr := mt.NewTimer(d)
		return mtmr.C, mtmr.Stop
	})
	// Let any timer that is still waiting finish when the test does.
	t.Cleanup(func() {
		mt.Advance(time.Hour)
	})
	return tmr, mt, called
}

func TestExpire(t *testing.T) {
	tmr, mt, called := newTimer(t, 5*ms)
	assert.Equal(t, timer.Running, tmr.State())
	assert.Equal(t, mt.Now().Add(5*ms), tmr.End())

	mt.Tick(4)
	assert.Equal(t, timer.Running, tmr.State())

	mt.Tick(1)
	recv(t, called)
	assert.Equal(t, timer.Done, tmr.State())
}

func TestUpdate(t *testing.T) {
	tmr, mt, called := newTimer(t, 5*ms)

	// Moving the end later replaces the original timer.
	end := mt.Now().Add(10 * ms)
	tmr.Update(end)
	assert.Equal(t, end, tmr.End())
	assert.Equal(t, 1, mt.Timers())

	mt.Tick(9)
	assert.Equal(t, timer.Running, tmr.State())

	mt.Tick(1)
	recv(t, called)
	assert.Equal(t, timer.Done, tmr.State())

	// Once the Timer is Done, Update does nothing.
	tmr.Update(mt.Now().Add(2 * ms))
	assert.Equal(t, timer.Done, tmr.State())
	assert.Equal(t, end, tmr.End())
	assert.Equal(t, 0, mt.Timers())
	mt.Tick(2)
	assert.Zero(t, len(called))
}

func TestUpdateDelay(t *testing.T) {
	tmr, mt, called := newTimer(t, 5*ms)

	// A Running Timer now expires 10ms from now instead of in 3ms.
	mt.Tick(2)
	tmr.UpdateDelay(10 * ms)
	assert.Equal(t, mt.Now().Add(10*ms), tmr.End())
	assert.Equal(t, 1, mt.Timers())

	mt.Tick(9)
	assert.Equal(t, timer.Running, tmr.State())

	// A Paused Timer stays Paused, with 3ms left for when it is resumed.
	tmr.Pause()
	tmr.UpdateDelay(3 * ms)
	assert.Equal(t, timer.Paused, tmr.State())
	assert.Equal(t, mt.Now().Add(3*ms), tmr.End())

	mt.Tick(5)
	assert.Zero(t, len(called))

	tmr.Resume()
	mt.Tick(3)
	recv(t, called)

	// Once the Timer is Done, UpdateDelay does nothing.
	end := tmr.End()
	tmr.UpdateDelay(2 * ms)
	assert.Equal(t, timer.Done, tmr.State())
	assert.Equal(t, end, tmr.End())
	assert.Equal(t, 0, mt.Timers())
}

func TestPauseResume(t *testing.T) {
	tmr, mt, called := newTimer(t, 5*ms)

	mt.Tick(2)
	tmr.Pause()
	assert.Equal(t, timer.Paused, tmr.State())
	assert.Equal(t, 0, mt.Timers())

	// While Paused, End keeps the 3ms that were left ahead of the clock.
	mt.Tick(10)
	assert.Equal(t, mt.Now().Add(3*ms), tmr.End())
	assert.Zero(t, len(called))

	tmr.Resume()
	assert.Equal(t, timer.Running, tmr.State())
	assert.Equal(t, mt.Now().Add(3*ms), tmr.End())

	mt.Tick(3)
	recv(t, called)
	assert.Equal(t, timer.Done, tmr.State())
}

func TestUpdatePaused(t *testing.T) {
	tmr, mt, called := newTimer(t, 5*ms)
	tmr.Pause()

	tmr.Update(mt.Now().Add(2 * ms))
	assert.Equal(t, timer.Running, tmr.State())

	mt.Tick(2)
	recv(t, called)
}

// Pause only applies to a Running Timer and Resume to a Paused one. Otherwise
// they do nothing.
func TestPauseResumeNoop(t *testing.T) {
	tmr, mt, called := newTimer(t, 5*ms)

	tmr.Resume()
	assert.Equal(t, timer.Running, tmr.State())
	assert.Equal(t, mt.Now().Add(5*ms), tmr.End())

	tmr.Pause()
	end := tmr.End()
	tmr.Pause()
	assert.Equal(t, end, tmr.End())

	tmr.Resume()
	mt.Tick(5)
	recv(t, called)

	tmr.Pause()
	tmr.Resume()
	assert.Equal(t, timer.Done, tmr.State())
}

func TestNilCallback(t *testing.T) {
	assert.PanicsWithValue(t, timer.ErrNilCallback, func() {
		timer.Delay(time.Second, nil)
	})
}

func TestUndefined(t *testing.T) {
	var tmr timer.Timer
	assert.Equal(t, timer.Undefined, tmr.State())
}

func TestNew(t *testing.T) {
	called := make(chan bool, 1)
	tmr := timer.New(time.Now().Add(ms), func() {
		called <- true
	})
	recv(t, called)
	assert.Equal(t, timer.Done, tmr.State())
}

func ExampleDelay() {
	done := make(chan bool)
	timer.Delay(time.Millisecond, func() {
		fmt.Println("the dishwasher is done")
		done <- true
	})
	<-done

	// Output: the dishwasher is done
}
