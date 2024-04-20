package cli

import "github.com/adamcolton/luce/util/handler"

// ExitClose holds what a program knows about being asked to exit or close, and
// what to do about it. Exit ends a client, and Close ends the service that the
// client is a client of.
type ExitClose struct {
	// Exit and Close are set when the client has been asked to exit or close.
	Exit, Close bool
	// CanExit and CanClose say if the exit and close commands are offered.
	CanExit, CanClose bool
	// RunExit and RunClose make the handler for the command call OnExit or
	// OnClose when the request is handled, rather than leaving it to a Runner.
	RunExit, RunClose bool
	// OnExit is called when a Runner exits, and OnClose after it if Close was
	// requested.
	OnExit, OnClose func()
}

// NewExitClose creates an ExitClose. CanExit is true if onExit is not nil, and
// CanClose is true if onClose is not.
func NewExitClose(onExit, onClose func()) *ExitClose {
	return &ExitClose{
		CanExit:  onExit != nil,
		CanClose: onClose != nil,
		OnExit:   onExit,
		OnClose:  onClose,
	}
}

// Commands creates the ExitCloseHandler for the exit and close commands, for a
// struct to embed.
func (ec *ExitClose) Commands() *ExitCloseHandler {
	return &ExitCloseHandler{
		ExitClose: ec,
		CloseDesc: "Close the server",
		ExitDesc:  "Exit the client",
	}
}

// EC returns the ExitClose, so that a struct that embeds one is a Commander.
func (ec *ExitClose) EC() *ExitClose {
	return ec
}

// ExitCloseHandler has the handlers and the usage for the exit and close
// commands, which are named for the methods. It is meant to be embedded in the
// struct that a handler.MethodsRegistrar reads. The ExitClose must be set. The
// descriptions are the usage of the commands.
type ExitCloseHandler struct {
	*ExitClose
	CloseDesc, ExitDesc string
}

// CloseReq is the request for the close command.
type CloseReq struct{}

// CloseResp is the response to the close command. A Runner that handles it
// closes and exits.
type CloseResp struct{}

// CloseHandler answers the close command. If RunClose is set it calls OnClose,
// if there is one.
func (ech *ExitCloseHandler) CloseHandler(e *CloseReq) *CloseResp {
	if ech.RunClose && ech.ExitClose.OnClose != nil {
		ech.ExitClose.OnClose()
	}
	return &CloseResp{}
}

// CloseUsage describes the close command. It is disabled if CanClose is false.
func (ech *ExitCloseHandler) CloseUsage() *handler.CommandDetails {
	return &handler.CommandDetails{
		Usage:    ech.CloseDesc,
		Disabled: !ech.CanClose,
	}
}

// ExitReq is the request for the exit command.
type ExitReq struct{}

// ExitResp is the response to the exit command. A Runner that handles it exits.
type ExitResp struct{}

// ExitHandler answers the exit command. If RunExit is set it calls OnExit, if
// there is one.
func (ech *ExitCloseHandler) ExitHandler(e *ExitReq) *ExitResp {
	if ech.RunExit && ech.ExitClose.OnExit != nil {
		ech.ExitClose.OnExit()
	}
	return &ExitResp{}
}

// ExitUsage describes the exit command. It is disabled if CanExit is false.
func (ech *ExitCloseHandler) ExitUsage() *handler.CommandDetails {
	return &handler.CommandDetails{
		Usage:    ech.ExitDesc,
		Disabled: !ech.CanExit,
	}
}
