package channel_test

import (
	"testing"
	"time"

	"github.com/adamcolton/luce/ds/channel"
	"github.com/adamcolton/luce/util/timeout"
	"github.com/stretchr/testify/assert"
)

func TestClose(t *testing.T) {
	c := channel.NewClose()

	done := make(chan bool)
	go func() {
		select {
		case <-c.OnClose:
			t.Error("c.OnClose triggered too soon")
		case <-time.After(time.Microsecond * 10):
		}
		done <- true
	}()

	timeout.Must(5, done)
	assert.False(t, c.Closed())

	go func() {
		select {
		case <-c.OnClose:
		case <-time.After(time.Millisecond):
			t.Error("c.OnClose did not trigger")
		}
		done <- true
	}()

	assert.True(t, c.Close())
	assert.False(t, c.Close())
	timeout.Must(5, done)

}
