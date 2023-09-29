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

// Listener combines a Receiver and a ListenerSwitcher to take in data from a
// bus, convert the data to an interface value and multiplex them out to the
// correct handlers.
type Listener interface {
	ListenerSwitcher
	// RegisterHandlers registers each handler with the ListenerSwitcher and the
	// type of its argument with the Receiver.
	RegisterHandlers(handler ...any) error
	// RegisterType tells the Receiver about a type of value it may receive.
	RegisterType(zeroValue any) error
}

// Sender handles the operations to place a message on a bus. For instance, it
// may contain the logic to serialize the message.
type Sender interface {
	// Send places the message on the bus.
	Send(msg any) error
}

// MultiSender will send a message to multiple busses at once. This can reduce
// duplication of work. For instance, if a message needs to be serialized, it
// will only be serialized once.
type MultiSender interface {
	// Send sends the message to the busses added with the ids. If no ids are
	// given it is sent to all of them, and an id that was not added is skipped.
	Send(msg any, ids ...string) error
	// Add adds a bus and names it key. It returns an error if 'to' is not a kind
	// of bus that the MultiSender accepts.
	Add(key string, to any) error
	// Delete removes the bus that was added as key.
	Delete(key string)
}
