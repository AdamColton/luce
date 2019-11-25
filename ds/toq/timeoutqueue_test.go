package toq_test

import (
	"testing"
	"time"

	"github.com/adamcolton/luce/ds/toq"
	"github.com/adamcolton/luce/util/timeout"
	"github.com/stretchr/testify/assert"
)

// The tests use timeout.MockTime as the queue's clock, so they never wait on
// real time. The clock only moves when a test calls Tick. The queue runs a Go
// routine that sleeps until the next timeout, and a test calls BlockUntil to
// wait for it to be asleep before moving the clock.

// newQueue creates a TimeoutQueue that uses a MockTime as its clock. Each tick
// of the clock is a millisecond.
func newQueue(t *testing.T, d time.Duration, capacity int) (*toq.TimeoutQueue, *timeout.MockTime) {
	mt := timeout.NewMockTime()
	tq := toq.New(d, capacity)
	tq.SetClock(mt.Now, mt.Sleep)
	// Wake any Go routine that is still sleeping when the test finishes.
	t.Cleanup(func() {
		mt.Advance(time.Hour)
	})
	return tq, mt
}

func getAction(ch chan<- int, i int) toq.TimeoutAction {
	return func() {
		ch <- i
	}
}

// recv gets a value from ch. The real time limit is only reached when a test
// fails.
func recv(t *testing.T, ch <-chan int) int {
	t.Helper()
	var i int
	assert.NoError(t, timeout.After(1000, func() {
		i = <-ch
	}))
	return i
}

func TestTimeoutQueue(t *testing.T) {
	d := time.Millisecond * 5
	tq, mt := newQueue(t, d, 10)

	assert.Equal(t, d, tq.Timeout())

	ch := make(chan int)
	tq.Add(getAction(ch, 1))
	mt.BlockUntil(1)
	mt.Tick(4)
	assert.Equal(t, 1, mt.Sleeping(), "the action should not be called early")
	mt.Tick(1)
	assert.Equal(t, 1, recv(t, ch))

	token1 := tq.Add(func() {
		t.Error("This should be canceled")
	})
	token2 := tq.Add(getAction(ch, 2))
	assert.True(t, token1.Cancel())
	assert.False(t, token1.Cancel())
	mt.BlockUntil(1)
	mt.Tick(5)
	assert.Equal(t, 2, recv(t, ch))
	assert.False(t, token2.Cancel())
}

func TestDecreaseSetTimeout(t *testing.T) {
	tq, mt := newQueue(t, time.Millisecond*100, 10)
	ch := make(chan int)

	tq.Add(getAction(ch, 0))
	tq.Add(getAction(ch, 1))
	tq.Add(getAction(ch, 2))
	mt.BlockUntil(1)

	// Nothing is due yet.
	mt.Tick(5)
	assert.Equal(t, 1, mt.Sleeping())

	// Everything in the queue is now overdue, so it should drain.
	tq.SetTimeout(time.Millisecond * 4)
	assert.Equal(t, time.Millisecond*4, tq.Timeout())
	// Cannot guarantee the order that the values will come through
	var got [3]bool
	got[recv(t, ch)] = true
	got[recv(t, ch)] = true
	got[recv(t, ch)] = true
	assert.Equal(t, [3]bool{true, true, true}, got)

	// The Go routine that was sleeping for the old timeout finds that another
	// has taken over.
	mt.Tick(100)
	assert.Zero(t, mt.Sleeping())

	tq.Add(getAction(ch, 4))
	mt.BlockUntil(1)
	mt.Tick(4)
	assert.Equal(t, 4, recv(t, ch))
}

func TestIncreaseSetTimeout(t *testing.T) {
	tq, mt := newQueue(t, time.Millisecond*10, 10)
	ch := make(chan int)

	tq.Add(getAction(ch, 1))
	mt.BlockUntil(1)
	mt.Tick(1)
	tq.Add(getAction(ch, 2))
	mt.Tick(1)
	tq.Add(getAction(ch, 3))

	// The three actions were added at 0, 1 and 2, so they are now due at 20,
	// 21 and 22.
	tq.SetTimeout(time.Millisecond * 20)

	// The Go routine wakes when the original timeout is up, finds that nothing
	// is due and goes back to sleep.
	mt.Tick(8)
	assert.Zero(t, mt.Sleeping())
	mt.BlockUntil(1)

	mt.Tick(10)
	assert.Equal(t, 1, recv(t, ch))
	mt.BlockUntil(1)
	mt.Tick(1)
	assert.Equal(t, 2, recv(t, ch))
	mt.BlockUntil(1)
	mt.Tick(1)
	assert.Equal(t, 3, recv(t, ch))
}

func TestReset(t *testing.T) {
	tq, mt := newQueue(t, time.Millisecond*20, 2)
	ch := make(chan int)

	// The actions are added at 0, 1 and 2, so they are due at 20, 21 and 22.
	tq.Add(getAction(ch, 1))
	mt.BlockUntil(1)
	mt.Tick(1)
	token := tq.Add(getAction(ch, 2))
	mt.Tick(1)
	tq.Add(getAction(ch, 3))
	mt.Tick(1)

	// Resetting at 3 makes the second action due at 23, after the third.
	assert.True(t, token.Reset())

	mt.Tick(17)
	assert.Equal(t, 1, recv(t, ch))
	mt.BlockUntil(1)
	mt.Tick(2)
	assert.Equal(t, 3, recv(t, ch))
	mt.BlockUntil(1)
	mt.Tick(1)
	assert.Equal(t, 2, recv(t, ch))

	// The action has run, so the token no longer does anything.
	assert.False(t, token.Reset())
	assert.False(t, token.Cancel())

	// The token's slot in the queue is reused for the next action, but the old
	// token can't touch it.
	tq.Add(getAction(ch, 4))
	assert.False(t, token.Reset())
	assert.False(t, token.Cancel())
	mt.BlockUntil(1)
	mt.Tick(20)
	assert.Equal(t, 4, recv(t, ch))

	// A token that was canceled can't be reset.
	token = tq.Add(getAction(ch, 5))
	assert.True(t, token.Cancel())
	assert.False(t, token.Reset())
}

func TestFlush(t *testing.T) {
	tq, mt := newQueue(t, time.Millisecond*5, 10)
	ch := make(chan int, 3)

	tq.Add(getAction(ch, 1))
	tq.Add(getAction(ch, 2))
	tq.Add(getAction(ch, 3))

	// Flush calls the actions itself, in order, before it returns.
	tq.Flush()
	assert.Equal(t, 1, recv(t, ch))
	assert.Equal(t, 2, recv(t, ch))
	assert.Equal(t, 3, recv(t, ch))

	// The queue is empty, so when the timeout is up nothing else is called.
	mt.BlockUntil(1)
	mt.Tick(5)
	assert.Zero(t, mt.Sleeping())
	assert.Zero(t, len(ch))

	// The queue still works after a Flush.
	tq.Add(getAction(ch, 4))
	mt.BlockUntil(1)
	mt.Tick(5)
	assert.Equal(t, 4, recv(t, ch))
}

func TestFlushMany(t *testing.T) {
	counts := make([]int, 5000)
	// The capacity is too small, so the queue has to grow.
	tq, _ := newQueue(t, time.Millisecond*5, 1)
	fn := func(i int) toq.TimeoutAction {
		return func() {
			counts[i]++
		}
	}
	for i := range counts {
		tq.Add(fn(i))
	}

	tq.Flush()
	for _, c := range counts {
		assert.Equal(t, 1, c)
	}
}

func TestFlushWhileFlushing(t *testing.T) {
	tq, _ := newQueue(t, time.Millisecond*5, 10)
	started := make(chan struct{})
	release := make(chan struct{})
	flushed := make(chan struct{})

	var ran int
	tq.Add(func() {
		close(started)
		<-release
		ran++
	})
	tq.Add(func() {
		ran++
	})

	go func() {
		tq.Flush()
		close(flushed)
	}()
	<-started

	// The first Flush is stuck in the first action, so this one returns
	// without doing anything.
	tq.Flush()
	assert.Zero(t, ran)

	close(release)
	<-flushed
	assert.Equal(t, 2, ran)
}
