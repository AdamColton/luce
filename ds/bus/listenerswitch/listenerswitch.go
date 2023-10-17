package listenerswitch

import (
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/handler"
)

const (
	ErrRegisterMuxHandler = lerr.Str("handler argument to RegisterMuxHandler requires a func or a channel")
)

// ListenerSwitch takes objects off a channel and passes them into the handlers
// for that type.
type ListenerSwitch struct {
	in <-chan any
	*handler.Switch
	lerr.ErrHandler
}

// New creates a ListenerSwitch for the in bus channel.
func New(size int, in <-chan any, errHandler any, handlers ...any) (*ListenerSwitch, error) {
	ls := &ListenerSwitch{
		Switch: handler.NewSwitch(size),
		in:     in,
	}
	var err error
	ls.ErrHandler, err = lerr.HandlerFunc(errHandler)
	if err != nil {
		return nil, err
	}
	for _, h := range handlers {
		if err := ls.RegisterInterface(h); err != nil {
			return nil, err
		}
	}
	return ls, nil
}

func (ls *ListenerSwitch) Handle(i any) (any, error) {
	out, err := ls.Switch.Handle(i)
	ls.ErrHandler.Handle(err)
	return out, err
}

func (ls *ListenerSwitch) SetErrorHandler(i any) (err error) {
	ls.ErrHandler, err = lerr.HandlerFunc(i)
	return
}

// Run the ListenerSwitch. It will stop running when the channel is closed.
func (ls *ListenerSwitch) Run() {
	for i := range ls.in {
		_, err := ls.Switch.Handle(i)
		ls.ErrHandler.Handle(err)
	}
}

// SetIn sets the interface channel the ListerMux is listening on.
func (ls *ListenerSwitch) SetIn(in <-chan interface{}) {
	ls.in = in
}
