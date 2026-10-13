package ldate

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUpdateNow(t *testing.T) {
	defer setNow(Now())

	var got []Date
	cancel := DateListener(func(d Date) { got = append(got, d) })
	defer cancel.Cancel()
	got = got[:0]

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
	assert.Equal(t, []Date{New(2024, 2, 28), New(2024, 2, 29)}, got)

	setNow(New(2024, 2, 29))
	assert.Len(t, got, 2, "no call when the date does not change")
}
