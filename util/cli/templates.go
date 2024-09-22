package cli

import (
	"github.com/adamcolton/luce/util/handler"
)

// HTMLTemplateLoadHandler has the handler for a command, loadHTMLTemplate, that
// calls Action, which is meant to reload the HTML templates. It is meant to be
// embedded in the struct that a handler.MethodsRegistrar reads.
type HTMLTemplateLoadHandler struct {
	Action  func()
	Details handler.CommandDetails
}

// NewHTMLTemplateLoadHandler creates an HTMLTemplateLoadHandler that calls
// action, which must not be nil. Its usage is "Reload HTML Templates".
func NewHTMLTemplateLoadHandler(action func()) *HTMLTemplateLoadHandler {
	return &HTMLTemplateLoadHandler{
		Action: action,
		Details: handler.CommandDetails{
			Usage: "Reload HTML Templates",
		},
	}
}

// HTMLTemplateLoadReq is the request for the loadHTMLTemplate command.
type HTMLTemplateLoadReq struct{}

// HTMLTemplateLoadResp is the response to the loadHTMLTemplate command.
type HTMLTemplateLoadResp struct{}

// LoadHTMLTemplateHandler answers the command by calling Action.
func (ldr *HTMLTemplateLoadHandler) LoadHTMLTemplateHandler(req *HTMLTemplateLoadReq) *HTMLTemplateLoadResp {
	ldr.Action()
	return &HTMLTemplateLoadResp{}
}

// LoadHTMLTemplateUsage describes the command with Details.
func (ldr *HTMLTemplateLoadHandler) LoadHTMLTemplateUsage() *handler.CommandDetails {
	return &(ldr.Details)
}
