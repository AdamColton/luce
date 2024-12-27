package ldate

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUpdateNow(t *testing.T) {
	defer setNow(Now())

	late := time.Date(2024, time.February, 28, 23, 59, 59, 999_000_000, time.UTC)
	morning := time.Date(2024, time.February, 29, 9, 0, 0, 0, time.UTC)
	stop := make(chan struct{})
	calls := 0
	clock := func() time.Time {
		calls++
		if calls == 1 {
			return late
		}
		close(stop)
		return morning
	}

	updateNow(clock, stop)
	assert.Equal(t, 2, calls, "it waited for the next day, then stopped")
	assert.Equal(t, New(2024, 2, 29), Now())
}
