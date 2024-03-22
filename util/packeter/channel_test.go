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

	// PackPipe closes its output when its input is closed.
	close(snd)
	timeout.Must(1000, func() { assert.Nil(t, <-rcv) })
}

func TestUnpackPipe(t *testing.T) {
	pre := prefix.New[uint32]()
	pip, snd, rcv := channel.NewPipe[[]byte](nil, nil)
	data := []byte("this is a test")
	go channel.Slice(pre.Pack(data), snd)
	go packeter.UnpackPipe(pre, pip, true)
	assert.Equal(t, data, <-rcv)

	// UnpackPipe closes its output when its input is closed.
	close(snd)
	timeout.Must(1000, func() { assert.Nil(t, <-rcv) })
}

// When cls is false, the pipe does not close chs.Snd when chs.Rcv closes.
func TestPipesDoNotClose(t *testing.T) {
	pre := prefix.New[uint32]()

	pip, snd, rcv := channel.NewPipe[[]byte](nil, nil)
	done := make(chan struct{})
	go func() {
		packeter.PackPipe(pre, pip, false)
		close(done)
	}()
	close(snd)
	assert.NoError(t, timeout.After(1000, done))
	select {
	case <-rcv:
		t.Error("PackPipe closed chs.Snd")
	default:
	}

	pip, snd, rcv = channel.NewPipe[[]byte](nil, nil)
	done = make(chan struct{})
	go func() {
		packeter.UnpackPipe(pre, pip, false)
		close(done)
	}()
	close(snd)
	assert.NoError(t, timeout.After(1000, done))
	select {
	case <-rcv:
		t.Error("UnpackPipe closed chs.Snd")
	default:
	}
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

	timeout.Must(1000, func() {
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

// unpackFunc is an Unpacker that is not also a Packer.
type unpackFunc func([]byte) [][]byte

func (fn unpackFunc) Unpack(data []byte) [][]byte {
	return fn(data)
}

func TestRunPackOnly(t *testing.T) {
	pipeTx, _, rcv := channel.NewPipe[[]byte](nil, nil)
	// prefix.Packer is a Packer and not an Unpacker.
	pipeOut := packeter.Run(&prefix.Packer[uint32]{}, pipeTx)
	assert.Nil(t, pipeOut.Rcv)

	data := []byte("this is a test")
	timeout.Must(1000, func() {
		pipeOut.Snd <- data
		// the length prefix, then the data
		assert.Len(t, <-rcv, 4)
		assert.Equal(t, data, <-rcv)

		close(pipeOut.Snd)
		assert.Nil(t, <-rcv)
	})
}

func TestRunUnpackOnly(t *testing.T) {
	pipeTx, snd, _ := channel.NewPipe[[]byte](nil, nil)
	// One chunk can complete more than one message.
	twice := unpackFunc(func(data []byte) [][]byte {
		return [][]byte{data, data}
	})
	pipeOut := packeter.Run(twice, pipeTx)
	assert.Nil(t, pipeOut.Snd)

	data := []byte("this is a test")
	timeout.Must(1000, func() {
		snd <- data
		assert.Equal(t, data, <-pipeOut.Rcv)
		assert.Equal(t, data, <-pipeOut.Rcv)

		close(snd)
		assert.Nil(t, <-pipeOut.Rcv)
	})
}

func TestRunNothingToDo(t *testing.T) {
	pipeTx, _, _ := channel.NewPipe[[]byte](nil, nil)

	// Neither a Packer nor an Unpacker
	pipeOut := packeter.Run(struct{}{}, pipeTx)
	assert.Nil(t, pipeOut.Snd)
	assert.Nil(t, pipeOut.Rcv)

	// A Packer and an Unpacker, but no channels to connect them to
	pipeOut = packeter.Run(prefix.New[uint32](), channel.Pipe[[]byte]{})
	assert.Nil(t, pipeOut.Snd)
	assert.Nil(t, pipeOut.Rcv)
}
