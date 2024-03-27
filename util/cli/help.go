package cli

type Helper string

type HelpReq struct {
	Command []string
}

func (h *HelpReq) Init(input []string) {
	h.Command = input
}

type HelpResp struct {
	Command []string
}

func (Helper) HelpHandler(req *HelpReq) *HelpResp {
	return &HelpResp{
		Command: req.Command,
	}
}

func (h Helper) HelpUsage() string {
	return string(h)
}

func (r *Runner) HelpRespHandler(resp *HelpResp) {
	r.ShowCommands(resp.Command)
}
