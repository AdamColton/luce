package bus

import (
	"github.com/adamcolton/luce/util/handler"
)

// ListenerSwitcher takes values off an interface channel and multiplexes them out
// to the correct handlers for the given type.
type ListenerSwitcher interface {
	handler.Switcher
	// Run takes values off the channel set by SetIn and passes each to the
	// handler for its type. It returns when the channel is closed.
	Run()
	// SetIn sets the channel that Run takes values from.
	SetIn(<-chan any)
	// SetErrorHandler sets what is called with each error, either a func(error) or
	// a channel of errors, as lerr.HandlerFunc accepts.
	SetErrorHandler(any) error
}

// Receiver receives data from a bus, translates it to a value and retransmits
// the value on an interface channel. For example, it may receive data as a byte
// slice, deserialize to a value and retransmit the value.
type Receiver interface {
	// Run receives data until the bus it reads from is closed. Nothing is
	// received until it is running.
	Run()
	// RegisterType tells the Receiver about a type of value it may receive.
	// zeroValue is a value of that type.
	RegisterType(zeroValue any) error
	// SetOut sets the channel that the values are sent on.
	SetOut(out chan<- any)
	// SetErrorHandler sets what is called with each error, either a func(error) or
	// a channel of errors, as lerr.HandlerFunc accepts.
	SetErrorHandler(any) error
}
