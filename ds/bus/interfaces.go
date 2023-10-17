package bus

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
