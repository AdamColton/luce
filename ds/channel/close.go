package channel

import "sync/atomic"

// Closed is the element type of Close.OnClose. Nothing is sent on that channel,
// it is closed to send the signal.
type Closed struct{}

// Close is a threadsafe close signal. Receiving on OnClose blocks until Close
// is called, so it can be used in a select. Calling Close more than once does
// not panic.
type Close struct {
	closed uint32
	// OnClose is closed when Close is called.
	OnClose chan Closed
}

// NewClose creates an instance of Close.
func NewClose() *Close {
	return &Close{
		OnClose: make(chan Closed),
	}
}

// Close is threadsafe and can be called multiple times, though only the first
// call will actually close the channel. The returned bool indicates if this
// call caused the channel to close.
func (c *Close) Close() bool {
	didClose := atomic.SwapUint32(&c.closed, 1) == 0
	if didClose {
		close(c.OnClose)
	}
	return didClose
}

// Closed checks if Close has been called. It reads the flag without
// synchronization, so it is only reliable when it is ordered after Close by
// something else, such as a receive from OnClose.
func (c *Close) Closed() bool {
	return c.closed == 1
}
