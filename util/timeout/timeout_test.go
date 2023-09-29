package timeout_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/adamcolton/luce/util/timeout"
	"github.com/stretchr/testify/assert"
)

func TestChan(t *testing.T) {
	intCh := make(chan int)
	go func() {
		intCh <- 10
	}()

	err := timeout.After(2, intCh)
	assert.NoError(t, err)

	err = timeout.After(2, intCh)
	assert.Equal(t, timeout.ErrTimeout, err)

	err = timeout.After(2, chan<- int(intCh))
	assert.Equal(t, timeout.ErrTimeout, err)

	go func() {
		assert.Equal(t, 0, <-intCh)
	}()
	err = timeout.After(2, chan<- int(intCh))
	assert.NoError(t, err)

	errCh := make(chan error)
	go func() {
		errCh <- nil
	}()

	err = timeout.After(2, errCh)
	assert.NoError(t, err)

	err = timeout.After(2, errCh)
	assert.Equal(t, timeout.ErrTimeout, err)

	go func() {
		errCh <- errors.New("testing")
	}()
	err = timeout.After(2, errCh)
	assert.Equal(t, "testing", err.Error())
}

func TestFunc(t *testing.T) {
	err := timeout.After(2, func() {})
	assert.NoError(t, err)

	err = timeout.After(2, func() {
		time.Sleep(time.Millisecond * 5)
	})
	assert.Equal(t, timeout.ErrTimeout, err)

	err = timeout.After(2, func() error {
		return errors.New("testing")
	})
	assert.Equal(t, "testing", err.Error())
}

func TestErrors(t *testing.T) {
	err := timeout.After(10, 3.1415)
	assert.Equal(t, fmt.Sprintf(timeout.InvalidWaitMsg, "float64"), err.Error())
}
