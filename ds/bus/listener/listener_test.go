package listener_test

import (
	"testing"

	"github.com/adamcolton/luce/ds/bus/listener"
	"github.com/adamcolton/luce/util/timeout"
	"github.com/stretchr/testify/assert"
)

type mockReceiver struct {
	out     chan<- any
	running bool
}

func (mr *mockReceiver) Run() {
	mr.running = true
	for mr.running {
	}
}
func (mr *mockReceiver) RegisterType(zeroValue any) error {
	return nil
}
func (mr *mockReceiver) SetOut(out chan<- any) {
	mr.out = out
}
func (mr *mockReceiver) SetErrorHandler(any) error {
	return nil
}

func TestListener(t *testing.T) {
	r := &mockReceiver{}
	errCh := make(chan error)
	strCh := make(chan string)
	handler := func(str string) {
		strCh <- str
	}
	l, err := listener.New(10, r, errCh, handler)
	assert.NoError(t, err)

	go l.Run()

	timeout.After(5, func() {
		str := "test"
		r.out <- str
		assert.Equal(t, str, <-strCh)
	})

}
