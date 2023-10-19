package timeout_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/adamcolton/luce/util/timeout"
	"github.com/stretchr/testify/assert"
)

func TestMockTimeNow(t *testing.T) {
	mt := timeout.NewMockTime()
	start := mt.Now()

	mt.Tick(3)
	assert.Equal(t, start.Add(3*time.Millisecond), mt.Now())

	mt.TickDuration = time.Second
	mt.Tick(2)
	assert.Equal(t, start.Add(3*time.Millisecond+2*time.Second), mt.Now())

	mt.Advance(time.Hour)
	assert.Equal(t, start.Add(3*time.Millisecond+2*time.Second+time.Hour), mt.Now())
}

func TestMockTimeSleep(t *testing.T) {
	mt := timeout.NewMockTime()

	// Three Go routines sleep for 10, 20 and 30 ticks. They report on woke
	// when they finish.
	woke := make(chan int)
	for i := 1; i <= 3; i++ {
		go func() {
			mt.Sleep(time.Duration(i*10) * time.Millisecond)
			woke <- i
		}()
	}
	mt.BlockUntil(3)
	assert.Equal(t, 3, mt.Sleeping())

	mt.Tick(9)
	assert.Equal(t, 3, mt.Sleeping(), "nobody should wake before their time")

	mt.Tick(1)
	assert.Equal(t, 2, mt.Sleeping())
	assert.NoError(t, timeout.After(1000, func() {
		assert.Equal(t, 1, <-woke)
	}))

	// Two sleeps finish at once, in either order.
	mt.Advance(20 * time.Millisecond)
	assert.Zero(t, mt.Sleeping())
	assert.NoError(t, timeout.After(1000, func() {
		got := map[int]bool{
			<-woke: true,
			<-woke: true,
		}
		assert.True(t, got[2])
		assert.True(t, got[3])
	}))
}

func TestMockTimeSleepNotPositive(t *testing.T) {
	mt := timeout.NewMockTime()
	assert.NoError(t, timeout.After(1000, func() {
		mt.Sleep(0)
		mt.Sleep(-time.Second)
	}))
	assert.Zero(t, mt.Sleeping())
}

func TestMockTimeBlockUntilAlreadyMet(t *testing.T) {
	mt := timeout.NewMockTime()
	go mt.Sleep(time.Millisecond)
	mt.BlockUntil(1)
	// Already met, so this returns without waiting.
	assert.NoError(t, timeout.After(1000, func() {
		mt.BlockUntil(1)
		mt.BlockUntil(0)
	}))
	mt.Tick(1)
}

func TestMockTimeManySleepers(t *testing.T) {
	mt := timeout.NewMockTime()

	wg := sync.WaitGroup{}
	wg.Add(1000)
	for i := 0; i < 1000; i++ {
		go func() {
			mt.Sleep(time.Duration(i%10+1) * time.Millisecond)
			wg.Done()
		}()
	}
	mt.BlockUntil(1000)

	for i := 0; i < 10; i++ {
		mt.Tick(1)
		assert.Equal(t, 1000-100*(i+1), mt.Sleeping())
	}
	assert.NoError(t, timeout.After(5000, &wg))
}

// Sleepers that start while the clock is moving are not lost.
func TestMockTimeConcurrentTicks(t *testing.T) {
	mt := timeout.NewMockTime()

	wg := sync.WaitGroup{}
	wg.Add(1000)
	go func() {
		for i := 0; i < 1000; i++ {
			go func() {
				mt.Sleep(time.Millisecond * 100)
				wg.Done()
			}()
			time.Sleep(time.Microsecond)
		}
	}()

	done := atomic.Bool{}
	ticked := make(chan struct{})
	go func() {
		defer close(ticked)
		for !done.Load() {
			mt.Tick(1)
			time.Sleep(time.Microsecond * 100)
		}
	}()

	assert.NoError(t, timeout.After(10000, &wg))
	done.Store(true)
	<-ticked
}

func TestMockTimeTimer(t *testing.T) {
	mt := timeout.NewMockTime()
	start := mt.Now()
	timer := mt.NewTimer(10 * time.Millisecond)
	assert.Equal(t, 1, mt.Timers())

	mt.Tick(9)
	select {
	case <-timer.C:
		t.Error("the timer should not have expired yet")
	default:
	}
	assert.Equal(t, 1, mt.Timers())

	mt.Tick(1)
	assert.Zero(t, mt.Timers())
	select {
	case got := <-timer.C:
		assert.Equal(t, start.Add(10*time.Millisecond), got)
	default:
		t.Error("the timer should have expired")
	}

	// It already expired, so there is nothing to stop.
	assert.False(t, timer.Stop())
}

func TestMockTimeTimerStop(t *testing.T) {
	mt := timeout.NewMockTime()
	first := mt.NewTimer(10 * time.Millisecond)
	second := mt.NewTimer(20 * time.Millisecond)
	third := mt.NewTimer(30 * time.Millisecond)

	assert.True(t, second.Stop())
	assert.False(t, second.Stop())
	assert.Equal(t, 2, mt.Timers())

	mt.Advance(time.Hour)
	assert.Zero(t, mt.Timers())
	assert.Len(t, first.C, 1)
	assert.Len(t, second.C, 0, "a stopped timer does not expire")
	assert.Len(t, third.C, 1)
}

func TestMockTimeTimerNotPositive(t *testing.T) {
	mt := timeout.NewMockTime()
	timer := mt.NewTimer(0)
	assert.Zero(t, mt.Timers())
	assert.Len(t, timer.C, 1)
	assert.False(t, timer.Stop())
}

func TestMockTimeBlockUntilTimers(t *testing.T) {
	mt := timeout.NewMockTime()

	// Already met, so this returns without waiting.
	mt.BlockUntilTimers(0)

	created := make(chan struct{})
	go func() {
		mt.NewTimer(time.Millisecond)
		mt.NewTimer(time.Millisecond)
		close(created)
	}()
	mt.BlockUntilTimers(2)
	assert.Equal(t, 2, mt.Timers())
	<-created
}
