package dbltimer_test

import (
	"testing"
	"time"

	"github.com/adamcolton/luce/util/dbltimer"
	"github.com/adamcolton/luce/util/timeout"
	"github.com/stretchr/testify/assert"
)

// Most of the tests use timeout.MockTime as the clock, so they never wait on
// real time. The clock only moves when a test calls Tick. The DoubleTimer's Go
// routine creates the soft timer when it starts, so newDoubleTimer waits for
// both timers before it returns.

const ms = time.Millisecond

// newDoubleTimer creates a DoubleTimer that uses a MockTime as its clock. The
// callback sends on called, which has a buffer so that the callback never
// blocks.
func newDoubleTimer(t *testing.T, hard, soft time.Duration) (*dbltimer.DoubleTimer, *timeout.MockTime, chan bool) {
	mt := timeout.NewMockTime()
	called := make(chan bool, 1)
	dt := dbltimer.NewWithTimer(hard, soft, func() {
		called <- true
	}, func(d time.Duration) (<-chan time.Time, func() bool) {
		timer := mt.NewTimer(d)
		return timer.C, timer.Stop
	})
	mt.BlockUntilTimers(2)
	// Let any timer that is still waiting finish when the test does.
	t.Cleanup(func() {
		mt.Advance(time.Hour)
	})
	return dt, mt, called
}

// recv waits for the callback. The real time limit is only reached when a test
// fails.
func recv(t *testing.T, called <-chan bool) {
	t.Helper()
	assert.NoError(t, timeout.After(1000, called), "the callback was not called")
}

func TestSoftLimit(t *testing.T) {
	dt, mt, called := newDoubleTimer(t, 20*ms, 2*ms)
	start := mt.Now()

	mt.Tick(1)
	assert.False(t, dt.Done())

	mt.Tick(1)
	recv(t, called)
	assert.True(t, dt.Done())

	// Triggered by the soft timer, well before the hard timer
	assert.Equal(t, 2*ms, mt.Now().Sub(start))
}

func TestReset(t *testing.T) {
	dt, mt, called := newDoubleTimer(t, 20*ms, 3*ms)

	// Each Reset starts the 3ms soft timer again, so it never gets to expire.
	for i := 0; i < 3; i++ {
		mt.Tick(1)
		assert.True(t, dt.Reset())
	}

	// The last Reset was 3ms ago in 3 ticks, so 2 ticks is not enough.
	mt.Tick(2)
	assert.False(t, dt.Done())

	mt.Tick(1)
	recv(t, called)
	assert.True(t, dt.Done())
}

func TestHardLimit(t *testing.T) {
	dt, mt, called := newDoubleTimer(t, 20*ms, 3*ms)
	start := mt.Now()

	// Resetting the soft timer over and over can't get past the hard timer.
	for i := 0; i < 19; i++ {
		mt.Tick(1)
		assert.True(t, dt.Reset())
	}
	assert.False(t, dt.Done())

	mt.Tick(1)
	recv(t, called)
	assert.True(t, dt.Done())

	// Triggered by the hard timer
	assert.Equal(t, 20*ms, mt.Now().Sub(start))

	// Once it has expired, Reset fails.
	assert.False(t, dt.Reset())
}

func TestCancel(t *testing.T) {
	dt, mt, called := newDoubleTimer(t, 20*ms, 2*ms)

	for i := 0; i < 5; i++ {
		mt.Tick(1)
		assert.True(t, dt.Reset())
	}

	assert.True(t, dt.Cancel())
	assert.True(t, dt.Done())

	// It can only be canceled once and can't be reset.
	assert.False(t, dt.Cancel())
	assert.False(t, dt.Reset())

	// Every timer expires, but the callback is not called.
	mt.Tick(100)
	assert.Zero(t, len(called))
}

func TestCancelAfterExpire(t *testing.T) {
	dt, mt, called := newDoubleTimer(t, 5*ms, 3*ms)

	mt.Tick(3)
	recv(t, called)

	// The callback was called, so it is too late to cancel.
	assert.False(t, dt.Cancel())
	assert.False(t, dt.Reset())
}

func TestNew(t *testing.T) {
	called := make(chan bool, 1)
	dt := dbltimer.New(time.Minute, ms, func() {
		called <- true
	})
	recv(t, called)
	assert.True(t, dt.Done())
}
