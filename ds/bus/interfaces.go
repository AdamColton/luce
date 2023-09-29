package bus

import (
	"github.com/adamcolton/luce/util/handler"
)

// ListenerSwitcher takes values off an interface channel and multiplexes them out
// to the correct handlers for the given type.
type ListenerSwitcher interface {
	handler.Switcher
	Run()
	SetIn(<-chan any)
	SetErrorHandler(any) error
}

// Receiver receives data from a bus translates it to a value and retransmits
// the value on an interface bus. For example, it may receive data as a byte
// slice, deserialize to a value and retransmit the value.
type Receiver interface {
	Run()
	RegisterType(zeroValue any) error
	SetOut(out chan<- any)
	SetErrorHandler(any) error
}

// Listener combines a Receiver and a ListenerSwitcher to take in data from a
// bus, convert the data to an interface value and multiplex them out to the
// correct handlers.
type Listener interface {
	ListenerSwitcher
	RegisterHandlers(handler ...any) error
	RegisterType(zeroValue any) error
}

// Sender handles the operations to place a message on a bus. For instance, it
// may contain the logic to serialize the message.
type Sender interface {
	Send(msg any) error
}

// MultiSender will send a message to multiple busses at once. This can reduce
// duplication of work. For instance, if a message needs to be serialized, it
// will only be serialized once.
type MultiSender interface {
	Send(msg any, ids ...string) error
	Add(key string, to any) error
	Delete(key string)
}
