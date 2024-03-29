package toq_test

import (
	"testing"
	"time"

	"github.com/adamcolton/luce/ds/toq"
	"github.com/adamcolton/luce/util/timeout"
	"github.com/stretchr/testify/assert"
)

// Note: While it would be better to avoid time.Sleep in tests, when the
// TimeoutQueue was updated to call the TimeoutActions in Go routines, it was
// difficult to guarentee the execution order. Adding small delays doesn't
// slow down the test much, but keeps the execution order predictible.

// TODO: replace direct calls to time.Now and time.Sleep

func getAction(ch chan<- int, i int) func() {
	time.Sleep(time.Millisecond)
	return func() {
		ch <- i
	}
}

func TestTimeoutQueue(t *testing.T) {
	d := time.Millisecond * 5
	tq := toq.New(d, 10)

	assert.Equal(t, d, tq.Timeout())

	ch := make(chan int)
	tq.Add(getAction(ch, 1))
	assert.NoError(t, timeout.After(20, func() {
		assert.Equal(t, 1, <-ch)
	}))

	token1 := tq.Add(func() {
		t.Error("This should be canceled")
	})
	token2 := tq.Add(getAction(ch, 2))
	assert.True(t, token1.Cancel())
	assert.NoError(t, timeout.After(20, ch))
	assert.False(t, token2.Cancel())
}

func TestDecreaseSetTimeout(t *testing.T) {
	tq := toq.New(time.Millisecond*100, 10)
	ch := make(chan int)

	tq.Add(getAction(ch, 0))
	tq.Add(getAction(ch, 1))
	tq.Add(getAction(ch, 2))

	// make sure there is some delay
	select {
	case <-ch:
		t.Error("too soon")
	case <-time.After(time.Millisecond * 5):
	}

	// Should cause the entire queue to drain
	tq.SetTimeout(time.Millisecond * 4)
	assert.NoError(t, timeout.After(3, func() {
		// Cannot guarentee the order that the values will come through
		var expected [3]bool
		expected[<-ch] = true
		expected[<-ch] = true
		expected[<-ch] = true
		assert.True(t, expected[0])
		assert.True(t, expected[1])
		assert.True(t, expected[2])
	}))

	tq.Add(getAction(ch, 4))
	assert.NoError(t, timeout.After(6, func() {
		assert.Equal(t, 4, <-ch)
	}))
}

func TestIncreaseSetTimeout(t *testing.T) {
	tq := toq.New(time.Millisecond*10, 10)
	ch := make(chan int)

	tq.Add(getAction(ch, 1))
	time.Sleep(time.Millisecond)
	tq.Add(getAction(ch, 2))
	time.Sleep(time.Millisecond)
	tq.Add(getAction(ch, 3))

	tq.SetTimeout(time.Millisecond * 20)

	// make sure there is some delay
	select {
	case <-ch:
		t.Error("too soon")
	case <-time.After(time.Millisecond * 8):
	}

	assert.NoError(t, timeout.After(20, func() {
		assert.Equal(t, 1, <-ch)
	}))
}

func TestFlush(t *testing.T) {
	tq := toq.New(time.Millisecond*5, 10)
	ch := make(chan int)

	tq.Add(getAction(ch, 1))
	tq.Add(getAction(ch, 2))
	tq.Add(getAction(ch, 3))

	done := make(chan bool)
	go func() {
		assert.Equal(t, 1, <-ch)
		assert.Equal(t, 2, <-ch)
		assert.Equal(t, 3, <-ch)
		done <- true
	}()

	tq.Flush()
	assert.NoError(t, timeout.After(1, done))
	// make sure nothing else sends on ch
	assert.Error(t, timeout.After(5, ch))
}

func TestDoubleFlush(t *testing.T) {
	counts := make([]int, 5000)
	tq := toq.New(time.Millisecond*5, len(counts)+1)
	fn := func(i int) func() {
		return func() {
			counts[i]++
		}
	}
	for i := range counts {
		tq.Add(fn(i))
	}
	done := make(chan bool)
	tq.Add(func() {
		done <- true
	})

	// Force Flush to be called multiple times
	go tq.Flush()
	go tq.Flush()
	assert.NoError(t, timeout.After(1, done))

	for _, c := range counts {
		assert.Equal(t, 1, c)
	}

}
