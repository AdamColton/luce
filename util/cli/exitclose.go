package cli

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

type ExitCloseHandler struct {
	*ExitClose
	CloseDesc, ExitDesc string
}

type CloseReq struct{}

type CloseResp struct{}

func (ech *ExitCloseHandler) CloseHandler(e *CloseReq) *CloseResp {
	return &CloseResp{}
}

func (ech *ExitCloseHandler) CloseUsage() (string, bool) {
	return ech.CloseDesc, ech.CanClose
}

type ExitReq struct{}

type ExitResp struct{}

func (ech *ExitCloseHandler) ExitHandler(e *ExitReq) *ExitResp {
	return &ExitResp{}
}

func (ech *ExitCloseHandler) ExitUsage() (string, bool) {
	return ech.ExitDesc, ech.CanExit
}

func (r *Runner) ExitRespHandler(e *ExitResp) {
	r.Exit = true
}

func (r *Runner) CloseRespHandler(c *CloseResp) {
	r.Close = true
	r.Exit = true
}
