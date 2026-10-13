package lerr

import "time"

// RetryPolicy configures Retry: how many times to call the function and how
// long to wait between calls.
type RetryPolicy struct {
	N     int
	Delay time.Duration
}

// DefaultRetryPolicy retries 5 times with a 1 millisecond delay.
var DefaultRetryPolicy = RetryPolicy{N: 5, Delay: time.Millisecond}

// Retry calls fn until it returns a nil error or the policy's attempts are
// used up, waiting Delay between calls. It returns fn's last result, so a
// caller can still inspect the final error on failure.
func Retry[T any](r RetryPolicy, fn func() (T, error)) (t T, err error) {
	for range r.N {
		t, err = fn()
		if err == nil {
			return
		}
		time.Sleep(r.Delay)
	}
	return
}
