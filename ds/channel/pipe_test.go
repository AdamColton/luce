package channel_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/channel"
	"github.com/stretchr/testify/assert"
)

func TestPipe(t *testing.T) {
	p, snd, rcv := channel.NewPipe[int](nil, nil)
	go func() {
		snd <- 10
		p.Snd <- 20
	}()

	assert.Equal(t, 10, <-p.Rcv)
	assert.Equal(t, 20, <-rcv)

	p2, i2, o2 := channel.NewPipe(rcv, snd)
	assert.Nil(t, i2)
	assert.Nil(t, o2)

	go func() {
		p2.Snd <- 30
		p.Snd <- 40
	}()
	assert.Equal(t, 30, <-p.Rcv)
	assert.Equal(t, 40, <-p2.Rcv)
}
