package channel_test

import (
	"sync"
	"sync/atomic"
	"testing"

	"github.com/adamcolton/luce/ds/channel"
	"github.com/stretchr/testify/assert"
)

func TestClose(t *testing.T) {
	c := channel.NewClose()

	select {
	case <-c.OnClose:
		t.Error("c.OnClose triggered too soon")
	default:
	}
	assert.False(t, c.Closed())

	assert.True(t, c.Close())
	assert.False(t, c.Close())

	select {
	case <-c.OnClose:
	default:
		t.Error("c.OnClose did not trigger")
	}
	assert.True(t, c.Closed())
}

func TestCloseConcurrent(t *testing.T) {
	c := channel.NewClose()

	// Many Go routines close it at once. It does not panic and exactly one of
	// them is the one that closed it.
	var closers atomic.Int32
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			if c.Close() {
				closers.Add(1)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	assert.Equal(t, int32(1), closers.Load())
	assert.True(t, c.Closed())
}
