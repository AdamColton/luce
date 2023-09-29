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
