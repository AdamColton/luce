package channel

import "sync/atomic"

// Closed is a flag.
type Closed struct{}

// Close uses a channel to send a Closed signal. It is threadsafe and will not
// panic if Close is called multiple times.
type Close struct {
	closed uint32
	// OnClose will block until Close is closed
	OnClose chan Closed
}

// NewClose creates an instance of Close.
func NewClose() *Close {
	return &Close{
		OnClose: make(chan Closed),
	}
}

// Close is threadsafe and can be call multiple times, though only the first
// call will actually close the channel. The returned bool indicates if this
// call caused the channel to close.
func (c *Close) Close() bool {
	didClose := atomic.SwapUint32(&c.closed, 1) == 0
	if didClose {
		close(c.OnClose)
	}
	return didClose
}

// Closed checks if Close has been called.
func (c *Close) Closed() bool {
	return c.closed == 1
}
