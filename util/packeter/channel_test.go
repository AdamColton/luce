package packeter_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/channel"
	"github.com/adamcolton/luce/util/packeter"
	"github.com/adamcolton/luce/util/packeter/prefix"
	"github.com/adamcolton/luce/util/timeout"
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

func TestRun(t *testing.T) {
	pre := prefix.New[uint32]()
	pipeTx, snd, rcv := channel.NewPipe[[]byte](nil, nil)
	pipeOut := packeter.Run(pre, pipeTx)

	data := []byte("this is a test")
	go func() {
		for r := range rcv {
			snd <- r
		}
		close(snd)
	}()

	timeout.Must(10, func() {
		pipeOut.Snd <- data
		got := <-pipeOut.Rcv
		assert.Equal(t, data, got)

		// This confirms that all pipes close correctly
		close(pipeOut.Snd)
		assert.Nil(t, <-rcv)
		assert.Nil(t, <-pipeOut.Rcv)
		assert.Nil(t, <-pipeTx.Rcv)
	})

}
