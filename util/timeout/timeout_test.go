package timeout_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/adamcolton/luce/util/timeout"
	"github.com/stretchr/testify/assert"
)

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
