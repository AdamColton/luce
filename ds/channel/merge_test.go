package channel_test

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/adamcolton/luce/ds/channel"
	"github.com/adamcolton/luce/util/timeout"
	"github.com/stretchr/testify/assert"
)

// Most of the tests replace the timer that ends a Cycle with a fakeTimer, so
// the test decides when a Cycle ends and nothing depends on real time.

// fakeTimer stands in for a dbltimer.DoubleTimer.
type fakeTimer struct {
	callback func()
	resets   atomic.Int32
}

func (f *fakeTimer) Reset() bool {
	f.resets.Add(1)
	return true
}

// expire ends the Cycle, as the real timer does when it expires. It blocks until
// the Cycle receives the signal.
func (f *fakeTimer) expire() {
	f.callback()
}

// newMerge creates a Merge that uses fakeTimers. The fakeTimer for each Cycle is
// sent on timers. Each time the Merge is closed and waits for the timer, it
// sends on waiting.
func newMerge() (m *channel.Merge[int], snd chan<- []int, rcv <-chan []int, timers <-chan *fakeTimer, waiting <-chan struct{}) {
	m, snd, rcv = channel.NewMerge[int](nil, nil)
	m.SingleDelay = time.Millisecond

	tch := make(chan *fakeTimer, 10)
	m.SetTimer(func(hard, soft time.Duration, callback func()) interface{ Reset() bool } {
		f := &fakeTimer{callback: callback}
		tch <- f
		return f
	})

	wch := make(chan struct{}, 1)
	m.SetTimeout(func(d time.Duration, ch <-chan bool) (bool, error) {
		select {
		case wch <- struct{}{}:
		default:
		}
		return channel.Timeout(d, ch)
	})
	return m, snd, rcv, tch, wch
}

// recvSlice gets a slice from rcv. The real time limit is only reached when a
// test fails.
func recvSlice(t *testing.T, rcv <-chan []int) []int {
	t.Helper()
	var s []int
	assert.NoError(t, timeout.After(1000, func() {
		s = <-rcv
	}))
	return s
}

// waitFor waits for done to be closed. The real time limit is only reached when
// a test fails.
func waitFor(t *testing.T, done <-chan struct{}) {
	t.Helper()
	assert.NoError(t, timeout.After(1000, done))
}

func TestMerge(t *testing.T) {
	m, snd, rcv, timers, waiting := newMerge()
	done := make(chan struct{})
	go func() {
		m.Run()
		close(done)
	}()

	// Slices that arrive while the timer is running are combined into one.
	expected := []int{1, 2, 3, 4, 6, 7, 8, 9}
	snd <- []int{1, 2}
	timer := <-timers
	snd <- []int{3, 4}
	snd <- []int{6, 7, 8, 9}
	// The Cycle finishes with the last slice before it sees the timer, so this
	// also waits until the timer has been reset for the last slice.
	timer.expire()
	assert.Equal(t, int32(2), timer.resets.Load(), "the first slice starts the timer, each other one resets it")
	assert.Equal(t, expected, recvSlice(t, rcv))

	// The next slice starts a new Cycle.
	snd <- []int{10}
	timer = <-timers
	timer.expire()
	assert.Equal(t, []int{10}, recvSlice(t, rcv))

	// Closing the incoming channel right after sending still lets the data go
	// through.
	snd <- []int{1, 2, 3, 4, 6, 7, 8, 9}
	timer = <-timers
	close(snd)
	<-waiting
	timer.expire()
	assert.Equal(t, expected, recvSlice(t, rcv))
	waitFor(t, done)

	// The outgoing channel is closed once Run is done.
	assert.Nil(t, recvSlice(t, rcv))

	// Once Run has finished, running it again doesn't panic.
	again := make(chan struct{})
	go func() {
		m.Run()
		close(again)
	}()
	waitFor(t, again)

	// Nothing is sent after that, even when a Cycle is started.
	cycled := make(chan struct{})
	go func() {
		m.Cycle([]int{1})
		close(cycled)
	}()
	waitFor(t, cycled)
}

func TestNewMerge(t *testing.T) {
	m, _, _ := channel.NewMerge[int](nil, nil)
	assert.Equal(t, channel.MaxDelay, m.MaxDelay)
	assert.Equal(t, channel.SingleDelay, m.SingleDelay)
}

// This one uses the real timers, so the delays are generous to keep it from
// depending on how quickly the Go routines run.
func TestMergeRealTimers(t *testing.T) {
	m, snd, rcv := channel.NewMerge[int](nil, nil)
	m.SingleDelay = 100 * time.Millisecond
	m.MaxDelay = time.Second
	go m.Run()

	snd <- []int{1}
	snd <- []int{2}
	snd <- []int{3}
	assert.Equal(t, []int{1, 2, 3}, recvSlice(t, rcv))
	close(snd)
}
