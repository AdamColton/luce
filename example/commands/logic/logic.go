package logic

import (
	"github.com/adamcolton/luce/util/handler"
)

type HandlerObject struct {
	People         []string
	EmptyCounter   int
	Closer, Exiter bool
	Timeout        int
}

type PersonReq struct {
	Name string
	Age  int
}

type PersonResp struct {
	Name string
}

func (ho *HandlerObject) PersonHandler(p *PersonReq) *PersonResp {
	ho.People = append(ho.People, p.Name)
	return &PersonResp{
		Name: p.Name,
	}
}

func (*HandlerObject) PersonUsage() string {
	return "Create a Person record"
}

type EmptyReq struct{}

func (ho *HandlerObject) EmptyHandler(e *EmptyReq) int {
	ho.EmptyCounter++
	return ho.EmptyCounter
}

func (*HandlerObject) EmptyUsage() string {
	return "Demonstrate empty struct handling"
}

type ListReq struct{}

func (ho *HandlerObject) ListHandler(e *ListReq) []string {
	return ho.People
}

func (*HandlerObject) ListUsage() string {
	return "List all the People"
}

type ExitReq struct{}

type ExitResp struct{}

func (ho *HandlerObject) ExitHandler(e *ExitReq) *ExitResp {
	return &ExitResp{}
}

func (ho *HandlerObject) ExitUsage() (string, bool) {
	return "Exit the client", ho.Exiter
}

type HelpReq struct{}

type HelpResp struct{}

func (ho *HandlerObject) HelpHandler(e *HelpReq) *HelpResp {
	return &HelpResp{}
}

func (*HandlerObject) HelpUsage() string {
	return "List all command"
}

type CloseReq struct{}

type CloseResp struct{}

func (ho *HandlerObject) CloseHandler(e *CloseReq) *CloseResp {
	return &CloseResp{}
}

func (ho *HandlerObject) CloseUsage() (string, bool) {
	return "Close the server", ho.Closer
}

type TimeoutReq struct{}

func (ho *HandlerObject) TimeoutHandler(e *TimeoutReq) {
	// No response causes a timeout
}

type SetTimeoutReq struct {
	Timeout int
}

type SetTimeoutResp struct {
	Timeout int
}

func (ho *HandlerObject) SetTimeoutHandler(e *SetTimeoutReq) *SetTimeoutResp {
	ho.Timeout = e.Timeout
	return &SetTimeoutResp{e.Timeout}
}

func (ho *HandlerObject) Commands() *handler.Commands {
	c := handler.DefaultRegistrar.Commands(ho)
	c = append(c, handler.Command{
		Name:  "",
		Usage: "",
		Action: func() *HelpResp {
			return &HelpResp{}
		},
	})

	cmds, _ := handler.Cmds(c)
	return cmds
}
