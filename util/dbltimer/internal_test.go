package dbltimer

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
