package timer

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/adamcolton/luce/util/timeout"
	"github.com/stretchr/testify/assert"
)

// Cancel can win the race with a timer: the timer expires, but by the time its
// callback runs the DoubleTimer has already been canceled. That can't be
// timed reliably from outside the package, so this sets the lock directly.
func TestCallbackAfterCancel(t *testing.T) {
	mt := timeout.NewMockTime()
	called := false
	dt := newDoubleTimer(time.Second, time.Second, func() {
		called = true
	}, func(d time.Duration) (<-chan time.Time, func() bool) {
		timer := mt.NewTimer(d)
		return timer.C, timer.Stop
	})
	mt.BlockUntilTimers(2)
	t.Cleanup(func() {
		mt.Advance(time.Hour)
	})

	atomic.StoreUint32(&dt.lock, cancel)
	dt.callback(hard)
	assert.False(t, called)
}

// Update and Pause can win the race with the timer: it expires, but before its
// Go routine gets the lock the run has been stopped. That can't be timed
// reliably from outside the package, so this calls expire with a stopped run.
func TestExpireAfterStop(t *testing.T) {
	mt := timeout.NewMockTime()
	called := false
	tmr := newTimerWithClock(mt.Now().Add(time.Second), func() {
		called = true
	}, mt.Now, func(d time.Duration) (<-chan time.Time, func() bool) {
		mtmr := mt.NewTimer(d)
		return mtmr.C, mtmr.Stop
	})
	stop := tmr.stop
	tmr.Pause()

	tmr.expire(stop)
	assert.False(t, called)
	assert.Equal(t, Paused, tmr.State())
}
