package timeout_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/adamcolton/luce/util/timeout"
	"github.com/stretchr/testify/assert"
)

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

func TestErrors(t *testing.T) {
	err := timeout.After(10, 3.1415)
	assert.Equal(t, fmt.Sprintf(timeout.InvalidWaitMsg, "float64"), err.Error())
}
