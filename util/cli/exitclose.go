package cli

import "github.com/adamcolton/luce/util/handler"

type ExitClose struct {
	Exit, Close       bool
	CanExit, CanClose bool
	OnExit, OnClose   func()
}

func NewExitClose(onExit, onClose func()) *ExitClose {
	return &ExitClose{
		CanExit:  onExit != nil,
		CanClose: onClose != nil,
		OnExit:   onExit,
		OnClose:  onClose,
	}
}

func (ec *ExitClose) Commands() *ExitCloseHandler {
	return &ExitCloseHandler{
		ExitClose: ec,
		CloseDesc: "Close the server",
		ExitDesc:  "Exit the client",
	}
}

func (ec *ExitClose) EC() *ExitClose {
	return ec
}

type ExitCloseHandler struct {
	*ExitClose
	CloseDesc, ExitDesc string
}

type CloseReq struct{}

type CloseResp struct{}

func (ech *ExitCloseHandler) CloseHandler(e *CloseReq) *CloseResp {
	return &CloseResp{}
}

func (ech *ExitCloseHandler) CloseUsage() *handler.CommandDetails {
	return &handler.CommandDetails{
		Usage:    ech.CloseDesc,
		Disabled: !ech.CanClose,
	}
}

type ExitReq struct{}

type ExitResp struct{}

func (ech *ExitCloseHandler) ExitHandler(e *ExitReq) *ExitResp {
	return &ExitResp{}
}

func (ech *ExitCloseHandler) ExitUsage() *handler.CommandDetails {
	return &handler.CommandDetails{
		Usage:    ech.ExitDesc,
		Disabled: !ech.CanExit,
	}
}
