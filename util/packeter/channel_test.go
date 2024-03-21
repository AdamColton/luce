package packeter_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/channel"
	"github.com/adamcolton/luce/util/packeter"
	"github.com/adamcolton/luce/util/packeter/prefix"
	"github.com/stretchr/testify/assert"
)

func TestPackPipe(t *testing.T) {
	pre := prefix.New[uint32]()
	pip, snd, rcv := channel.NewPipe[[]byte](nil, nil)

	go packeter.PackPipe(pre, pip, true)
	data := []byte("this is a test")
	snd <- data
	pre.Unpack(<-rcv)
	got := pre.Unpack(<-rcv)
	assert.Len(t, got, 1)
	assert.Equal(t, data, got[0])
	close(snd)
}

func TestUnpackPipe(t *testing.T) {
	pre := prefix.New[uint32]()
	pip, snd, rcv := channel.NewPipe[[]byte](nil, nil)
	data := []byte("this is a test")
	go channel.Slice(pre.Pack(data), snd)
	go packeter.UnpackPipe(pre, pip, true)
	assert.Equal(t, data, <-rcv)
	close(snd)
}
