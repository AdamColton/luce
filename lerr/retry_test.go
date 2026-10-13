package lerr_test

import (
	"testing"
	"time"

	"github.com/adamcolton/luce/lerr"
	"github.com/stretchr/testify/assert"
)

func TestRetrySucceedsAfterFailures(t *testing.T) {
	attempts := 0
	flaky := func() (string, error) {
		attempts++
		if attempts < 3 {
			return "", lerr.Str("not ready")
		}
		return "ready", nil
	}

	got, err := lerr.Retry(lerr.RetryPolicy{N: 5, Delay: time.Millisecond}, flaky)
	assert.NoError(t, err)
	assert.Equal(t, "ready", got)
	assert.Equal(t, 3, attempts)
}

func TestRetryExhausted(t *testing.T) {
	alwaysFails := func() (string, error) {
		return "", lerr.Str("still not ready")
	}

	got, err := lerr.Retry(lerr.RetryPolicy{N: 3, Delay: time.Millisecond}, alwaysFails)
	assert.EqualError(t, err, "still not ready")
	assert.Equal(t, "", got)
}
