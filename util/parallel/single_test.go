package parallel_test

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/adamcolton/luce/util/parallel"
	"github.com/stretchr/testify/assert"
)

func TestSingle(t *testing.T) {
	calls := 0
	s := parallel.NewSingle(func() { calls++ })
	assert.False(t, s.Running())

	s.Run()
	assert.Equal(t, 1, calls)
	assert.False(t, s.Running())

	s.Run()
	assert.Equal(t, 2, calls)
}

func TestSingleConcurrent(t *testing.T) {
	var calls atomic.Int32
	started := make(chan bool)
	release := make(chan bool)
	var s *parallel.Single
	s = parallel.NewSingle(func() {
		calls.Add(1)
		started <- s.Running()
		<-release
	})

	done := make(chan struct{})
	go func() {
		s.Run()
		close(done)
	}()
	assert.True(t, <-started)

	// While fn is blocked, other calls return immediately without calling fn.
	wg := &sync.WaitGroup{}
	for range 10 {
		wg.Add(1)
		go func() {
			s.Run()
			wg.Done()
		}()
	}
	wg.Wait()
	assert.Equal(t, int32(1), calls.Load())

	close(release)
	<-done
	assert.False(t, s.Running())

	// Once finished, fn can be run again.
	release = make(chan bool)
	close(release)
	go func() { <-started }()
	s.Run()
	assert.Equal(t, int32(2), calls.Load())
}

func TestSinglePanic(t *testing.T) {
	calls := 0
	s := parallel.NewSingle(func() {
		calls++
		panic("test")
	})
	assert.Panics(t, s.Run)
	assert.True(t, s.Running())

	// fn will not be invoked again after a panic.
	s.Run()
	assert.Equal(t, 1, calls)
}
