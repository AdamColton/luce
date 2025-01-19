package dbltimer_test

import (
	"testing"
	"time"

	"github.com/adamcolton/luce/util/dbltimer"
	"github.com/adamcolton/luce/util/timeout"
	"github.com/stretchr/testify/assert"
)

// TODO: to really test the logic of race conditions I should create
// a var for time.NewTimer and create some internal tests

func TestSoftLimit(t *testing.T) {
	ms := time.Millisecond
	called := make(chan bool)
	callback := func() {
		called <- true
	}
	start := time.Now()
	dbltimer.New(20*ms, 2*ms, callback)

	to := timeout.After(50, called)
	assert.NoError(t, to)
	d := time.Since(start)

	// Triggered by soft limit not hard limit
	assert.True(t, d < ms*19)
}

func TestReset(t *testing.T) {
	ms := time.Millisecond
	called := make(chan bool)
	callback := func() {
		called <- true
	}
	start := time.Now()
	dt := dbltimer.New(20*ms, 3*ms, callback)

	for i := 0; i < 3; i++ {
		time.Sleep(ms)
		assert.True(t, dt.Reset())
	}

	to := timeout.After(50, called)
	assert.NoError(t, to)
	d := time.Since(start)

	// Triggered by soft limit not hard limit
	assert.True(t, d < ms*19)
	assert.True(t, dt.Done())
}

func TestHardLimit(t *testing.T) {
	ms := time.Millisecond
	called := make(chan bool)
	callback := func() {
		called <- true
	}
	start := time.Now()
	dt := dbltimer.New(20*ms, 3*ms, callback)

	resetFail := make(chan bool)
	go func() {
		for dt.Reset() {
			time.Sleep(ms)
		}
		resetFail <- true
	}()

	to := timeout.After(50, called)
	assert.NoError(t, to)
	d := time.Since(start)

	// Triggered by hard limit
	assert.True(t, d > ms*19)
	gotResetFail := timeout.After(100, resetFail)
	assert.NoError(t, gotResetFail)
}

func TestCancel(t *testing.T) {
	wrap := func() {
		ms := time.Millisecond
		called := make(chan bool)
		callback := func() {
			called <- true
		}

		dt := dbltimer.New(20*ms, 2*ms, callback)

		// spam reset during a cancel to check if it triggers a race condition
		spam := make(chan bool)
		count := 0
		go func() {
			for dt.Reset() {
				count++
			}
			spam <- true
		}()

		time.Sleep(ms)
		assert.True(t, dt.Cancel())

		to := timeout.After(30, called)
		assert.Equal(t, timeout.ErrTimeout, to)

		to = timeout.After(30, spam)
		assert.NoError(t, to)

		// make sure there was a reasonable degree of spam
		assert.Greater(t, count, 10)
	}

	err := timeout.After(100, wrap)
	assert.NoError(t, err)
}
