package cli

import (
	"github.com/adamcolton/luce/util/handler"
)

type HTMLTemplateLoadHandler struct {
	Action  func()
	Details handler.CommandDetails
}

func NewHTMLTemplateLoadHandler(action func()) *HTMLTemplateLoadHandler {
	return &HTMLTemplateLoadHandler{
		Action: action,
		Details: handler.CommandDetails{
			Usage: "Reload HTML Templates",
		},
	}
}

type HTMLTemplateLoadReq struct{}

type HTMLTemplateLoadResp struct{}

func (ldr *HTMLTemplateLoadHandler) LoadHTMLTemplateHandler(req *HTMLTemplateLoadReq) *HTMLTemplateLoadResp {
	ldr.Action()
	return &HTMLTemplateLoadResp{}
}

func (ldr *HTMLTemplateLoadHandler) LoadHTMLTemplateUsage() *handler.CommandDetails {
	return &(ldr.Details)
}
