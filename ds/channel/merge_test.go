package channel_test

import (
	"sync"
	"testing"
	"time"

	"github.com/adamcolton/luce/ds/channel"
	"github.com/adamcolton/luce/util/timeout"
	"github.com/stretchr/testify/assert"
)

func TestMerge(t *testing.T) {
	m, snd, rcv := channel.NewMerge[int](nil, nil)
	m.SingleDelay = 10000 * time.Microsecond
	m.MaxDelay = 100 * time.Millisecond
	expected := []int{1, 2, 3, 4, 6, 7, 8, 9}

	running := true
	go func() {
		m.Run()
		running = false
	}()

	// sending expected in 3 chunks with small delay is correct reassmbled into
	// expected
	go func() {
		snd <- expected[:2]
		time.Sleep(time.Microsecond * 10)
		snd <- expected[2:4]
		time.Sleep(time.Microsecond * 10)
		snd <- expected[4:]
	}()
	assert.Equal(t, expected, <-rcv)

	// send 2*MaxDelay/SingleDelay slices (accounting for ms/us = 1000)
	// sending a slice twice per SingleDelayUS keeps the Cycle open for
	// MaxDelayMS once, then the rest make it through on the second cycle
	// ideally this would only require
	toSend := int((m.MaxDelay * 3) / (m.SingleDelay))
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		d := m.SingleDelay / 2
		for i := 0; i < toSend; i++ {
			time.Sleep(d)
			snd <- []int{i}
		}
		wg.Done()
	}()
	timeout.Must(int(m.MaxDelay*5), func() {
		ln := len(<-rcv)
		assert.True(t, ln < toSend)
		assert.True(t, ln > 0)
		ln += len(<-rcv)
		assert.Equal(t, toSend, ln)
	})

	// confirm that closing right after sending still allows the data
	// to go through
	go func() {
		wg.Wait()
		snd <- expected
		close(snd)
	}()

	err := timeout.After(100, func() {
		assert.Equal(t, expected, <-rcv)
	})
	assert.NoError(t, err)
	assert.False(t, running)

	// confirm that calling m.Run again doesn't panic
	timeout.Must(1, m.Run)
}
