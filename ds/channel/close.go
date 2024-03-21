package channel

import "sync/atomic"

// Closed is a flag.
type Closed struct{}

type Close struct {
	closed uint32
	// OnClose will block until Close is closed
	OnClose chan Close
}

func NewClose() *Close {
	return &Close{
		OnClose: make(chan Close),
	}
}

func (c *Close) Close() bool {
	didClose := atomic.SwapUint32(&c.closed, 1) == 0
	if didClose {
		close(c.OnClose)
	}
	return didClose
}

func (c *Close) Closed() bool {
	return c.closed == 1
}
