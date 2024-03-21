package cli

import "github.com/adamcolton/luce/util/handler"

// Helper is the description of the help command, and has the handler for it. It
// is meant to be embedded in the struct that a handler.MethodsRegistrar reads.
type Helper string

// HelpReq is the request for the help command. Command is the path of the
// command to list the sub-commands of, or empty to list them all.
type HelpReq struct {
	Command []string
}

// Init takes the words after the command as the path.
func (h *HelpReq) Init(input []string) {
	h.Command = input
}

// HelpResp is the response to the help command. A Runner that handles it writes
// the commands.
type HelpResp struct {
	Command []string
}

// HelpHandler answers the help command.
func (Helper) HelpHandler(req *HelpReq) *HelpResp {
	return &HelpResp{
		Command: req.Command,
	}
}

// HelpUsage describes the help command with the Helper's text.
func (h Helper) HelpUsage() *handler.CommandDetails {
	return &handler.CommandDetails{
		Usage: string(h),
	}
}

// HelpRespHandler is the handler for a HelpResp. It calls ShowCommands.
func (r *Runner) HelpRespHandler(resp *HelpResp) {
	r.ShowCommands(resp.Command)
}
