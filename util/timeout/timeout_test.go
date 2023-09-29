package timeout_test

import (
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/adamcolton/luce/util/timeout"
	"github.com/stretchr/testify/assert"
)

func TestChan(t *testing.T) {
	intCh := make(chan int)
	go func() {
		intCh <- 10
	}()

	err := timeout.After(1000, intCh)
	assert.NoError(t, err)

	err = timeout.After(2, intCh)
	assert.Equal(t, timeout.ErrTimeout, err)

	err = timeout.After(2, chan<- int(intCh))
	assert.Equal(t, timeout.ErrTimeout, err)

	go func() {
		assert.Equal(t, 0, <-intCh)
	}()
	err = timeout.After(1000, chan<- int(intCh))
	assert.NoError(t, err)

	errCh := make(chan error)
	go func() {
		errCh <- nil
	}()

	err = timeout.After(1000, errCh)
	assert.NoError(t, err)

	err = timeout.After(2, errCh)
	assert.Equal(t, timeout.ErrTimeout, err)

	go func() {
		errCh <- errors.New("testing")
	}()
	err = timeout.After(1000, errCh)
	assert.Equal(t, "testing", err.Error())
}

func TestFunc(t *testing.T) {
	err := timeout.After(1000, func() {})
	assert.NoError(t, err)

	// The function does not return until it is released, so it times out.
	release := make(chan struct{})
	err = timeout.After(2, func() {
		<-release
	})
	assert.Equal(t, timeout.ErrTimeout, err)
	close(release)

	err = timeout.After(1000, func() error {
		return errors.New("testing")
	})
	assert.Equal(t, "testing", err.Error())
}

func TestWaitGroup(t *testing.T) {
	wg := &sync.WaitGroup{}

	err := timeout.After(1000, wg)
	assert.NoError(t, err)

	wg.Add(1)
	go wg.Done()
	err = timeout.After(1000, wg)
	assert.NoError(t, err)

	// Nothing calls Done, so it times out.
	wg.Add(1)
	err = timeout.After(4, wg)
	assert.Equal(t, timeout.ErrTimeout, err)
	wg.Done()
}

func TestErrors(t *testing.T) {
	err := timeout.After(10, 3.1415)
	assert.Equal(t, fmt.Sprintf(timeout.InvalidWaitMsg, "float64"), err.Error())
}
