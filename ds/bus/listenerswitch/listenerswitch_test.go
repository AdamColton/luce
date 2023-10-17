package listenerswitch_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/bus/listenerswitch"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/handler"
	"github.com/adamcolton/luce/util/timeout"
	"github.com/stretchr/testify/assert"
)

func TestListenerMux(t *testing.T) {
	strCh := make(chan string)
	ch := make(chan any)
	done := make(chan bool)

	handler := func(s string) {
		strCh <- s
	}
	ls, err := listenerswitch.New(10, ch, nil, handler)
	assert.NoError(t, err)

	runner := func() {
		ls.Run()
		done <- true
	}
	testFn := func() {
		str := "test"
		ch <- str
		assert.Equal(t, "test", <-strCh)
		close(ch)
		assert.True(t, <-done)
	}

	go runner()
	assert.NoError(t, timeout.After(5, testFn))

	ch = make(chan any)
	ls.SetIn(ch)
	go runner()
	assert.NoError(t, timeout.After(5, testFn))

}

func TestListenerMuxErr(t *testing.T) {
	ls, err := listenerswitch.New(10, nil, nil, 123)
	assert.Nil(t, ls)
	assert.Equal(t, handler.ErrRegisterInterface, err)

	ls, err = listenerswitch.New(10, nil, "abc", 123)
	assert.Nil(t, ls)
	assert.Equal(t, lerr.ErrHandlerFunc, err)
}
