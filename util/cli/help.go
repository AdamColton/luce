package cli

type Helper string

type HelpReq struct{}

type HelpResp struct{}

func (Helper) HelpHandler(e *HelpReq) *HelpResp {
	return &HelpResp{}
}

func (h Helper) HelpUsage() string {
	return string(h)
}

func (r *Runner) HelpRespHandler(h *HelpResp) {
	r.ShowCommands()
}
