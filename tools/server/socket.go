package server

import (
	"fmt"
	"strings"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/cli"
	"github.com/adamcolton/luce/util/handler"
)

// RunSocket for the admin interface. This is not invoked by ListenAndServe
// and needs to be run seperately.
// func (s *Server) RunSocket() {
// 	sck := unixsocket.New(s.Socket, func(conn net.Conn) {
// 		pipe := unixsocket.ConnPipe(conn)
// 		w := channel.Writer{pipe.Snd}

// 		ctx := cli.NewContext(w, pipe.Rcv, nil)
// 		s.Cli(ctx, func() {
// 			conn.Close()
// 		})
// 	})

// 	sck.Run()
// }

// func (s *Server) Cli(ctx cli.Context, onExit func()) {
// 	onClose := func() {
// 		s.Close()
// 	}
// 	ec := cli.NewExitClose(onExit, onClose)
// 	c := &cliHandlers{
// 		Server:           s,
// 		ExitCloseHandler: ec.Commands(),
// 	}

// 	r := cli.NewRunner(c, ctx)
// 	r.Prompt = "> "
// 	r.StartMessage = "Welcome to the luce server\nenter 'help' for more\n"
// 	r.Run()
// }

func (s *Server) coreCommander(ec *cli.ExitClose) cli.Commander {
	return &cliHandlers{
		Server:           s,
		ExitCloseHandler: ec.Commands(),
	}
}

func (s *Server) RunStdIO() {
	s.coreserver.RunStdIO()
}

func (c *cliHandlers) Handlers(rnr *cli.Runner) []any {
	return []any{
		func(r *CreateUserResp) {
			if r.Error != nil {
				rnr.WriteStrings("Failed to create User: ", r.Error.Error())
			} else {
				rnr.WriteString("Created User")
			}
		},
		func(r ListUsersResp) {
			fmt.Fprintf(rnr, "  %s", strings.Join(r, "\n  "))
		},
		func(r *GroupResp) {
			if r.Error != nil {
				rnr.WriteStrings("Failed to create Group: ", r.Error.Error())
			} else {
				rnr.WriteString("Created Group")
			}
		},
		func(r ListGroupsResp) {
			fmt.Fprintf(rnr, "  %s", strings.Join(r, "\n  "))
		},
		func(r *UserGroupResp) {
			if r.Error != nil {
				rnr.WriteStrings("Failed to add User to Group: ", r.Error.Error())
			} else {
				rnr.WriteString("Added User to Group")
			}
		},
		func(r *SetPortResp) {
			rnr.WriteString("Port changed, server restarted")
		},
		func(r Settings) {
			fmt.Fprintf(rnr, "  AdminLockUserCreation %t", r.AdminLockUserCreation)
		},
		func(r *RoutesResp) {
			fmt.Fprint(rnr, string(*r))
		},
		func(r ListServicesResp) {
			fmt.Fprintf(rnr, "  %s", strings.Join(r, "\n  "))
		},
		func(r ListBashCommandsResp) {
			for i, bc := range r {
				fmt.Fprintf(rnr, "  %d: %s (running:%t auto:%t)\n", i, bc.Format, bc.running, bc.Auto)
			}
		},
		func(r RunBashCommandResp) {
			fmt.Fprintln(rnr, r.Msg)
		},
		func(r BashCommandOuputResp) {
			fmt.Fprintln(rnr, r.Msg)
		},
		rnr.ExitRespHandler,
		rnr.CloseRespHandler,
		rnr.HelpRespHandler,
	}
}

type cliHandlers struct {
	Server *Server
	*cli.ExitCloseHandler
	cli.Helper
}

type CreateUserReq struct {
	Name, Password string
}

type CreateUserResp struct {
	Error error
}

func (c *cliHandlers) UserHandler(req *CreateUserReq) *CreateUserResp {
	_, err := c.Server.Users.Create(req.Name, req.Password)
	return &CreateUserResp{Error: err}
}

func (*cliHandlers) UserUsage() *handler.CommandDetails {
	return &handler.CommandDetails{
		Usage: "Create a User",
		Alias: "u",
	}
}

type ListUsersReq struct{}
type ListUsersResp []string

func (c *cliHandlers) ListUsersHandler(req *ListUsersReq) ListUsersResp {
	return c.Server.Users.List()
}

func (*cliHandlers) ListUsersUsage() *handler.CommandDetails {
	return &handler.CommandDetails{
		Usage: "List User names",
		Alias: "lu",
	}
}

type GroupReq struct {
	Name string
}
type GroupResp struct {
	Error error
}

func (c *cliHandlers) GroupHandler(req *GroupReq) *GroupResp {
	_, err := c.Server.Users.Group(req.Name)
	return &GroupResp{Error: err}
}

func (*cliHandlers) GroupUsage() *handler.CommandDetails {
	return &handler.CommandDetails{
		Usage: "Create Group",
		Alias: "g",
	}
}

type ListGroupsReq struct{}
type ListGroupsResp []string

func (c *cliHandlers) ListGroupsHandler(req *ListGroupsReq) ListGroupsResp {
	return c.Server.Users.Groups()
}

func (*cliHandlers) ListGroupsUsage() *handler.CommandDetails {
	return &handler.CommandDetails{
		Usage: "List Groups",
		Alias: "lg",
	}
}

type UserGroupReq struct {
	User, Group string
}
type UserGroupResp struct {
	Error error
}

func (c *cliHandlers) UserGroupHandler(req *UserGroupReq) *UserGroupResp {

	g := c.Server.Users.HasGroup(req.Group)
	if g == nil {
		return &UserGroupResp{
			Error: lerr.Str("group not found"),
		}
	}

	u, err := c.Server.Users.GetByName(req.User)
	if err != nil {
		return &UserGroupResp{
			Error: err,
		}
	}
	if u == nil {
		return &UserGroupResp{
			Error: lerr.Str("user not found"),
		}
	}

	return &UserGroupResp{
		Error: g.AddUser(u),
	}
}

func (*cliHandlers) UserGroupUsage() *handler.CommandDetails {
	return &handler.CommandDetails{
		Usage: "Add User to Group",
		Alias: "ug",
	}
}

type SetPortReq struct {
	Port string
}
type SetPortResp struct{}

func (c *cliHandlers) SetPortHandler(req *SetPortReq) *SetPortResp {
	c.Server.Close()
	c.Server.coreserver.Addr = req.Port
	go c.Server.coreserver.ListenAndServe()
	return &SetPortResp{}
}

func (*cliHandlers) SetPortUsage() *handler.CommandDetails {
	return &handler.CommandDetails{
		Usage: "Set server port",
		Alias: "sp",
	}
}

type SettingsReq struct{}
type Settings struct {
	AdminLockUserCreation bool
}

func (c *cliHandlers) SettingsHandler(req *SettingsReq) Settings {
	return c.Server.Settings
}

func (*cliHandlers) SettingsUsage() *handler.CommandDetails {
	return &handler.CommandDetails{
		Usage: "Display server settings",
		Alias: "s",
	}
}

type AdminLockUserCreationReq struct {
	AdminLockUserCreation bool
}
type AdminLockUserCreationResp struct{}

func (c *cliHandlers) AdminLockUserCreationHandler(req *AdminLockUserCreationReq) *AdminLockUserCreationResp {
	c.Server.Settings.AdminLockUserCreation = req.AdminLockUserCreation
	return &AdminLockUserCreationResp{}
}

func (*cliHandlers) AdminLockUserCreationUsage() *handler.CommandDetails {
	return &handler.CommandDetails{
		Usage: "Restricts user creation to admins",
		Alias: "su",
	}
}

func (c *cliHandlers) Commands() *handler.Commands {
	cmds := handler.DefaultRegistrar.Commands(c)
	if x, ok := cmds.Get("exit"); ok {
		x.Alias = "q"
	}
	if q, ok := cmds.Get("close"); ok {
		q.Alias = "cls"
	}
	if h, ok := cmds.Get("help"); ok {
		h.Alias = "h"
	}
	cs := cmds.Vals(nil).Sort(handler.CmdNameLT)

	return lerr.Must(handler.Cmds(cs))
}

type RoutesReq struct {
}
type RoutesResp string

func (c *cliHandlers) RoutesHandler(req *RoutesReq) *RoutesResp {

	out := make(slice.Slice[string], 0, c.Server.serviceRoutes.Len())
	c.Server.serviceRoutes.Each(func(k string, route *serviceRoute, done *bool) {
		if route.active {
			out = append(out, k)
		}
	})
	out.Sort(slice.LT[string]())
	resp := RoutesResp(strings.Join(out, "\n"))
	return &resp
}

func (*cliHandlers) RoutesUsage() *handler.CommandDetails {
	return &handler.CommandDetails{
		Usage: "Show Routes",
		Alias: "sr",
	}
}

type ListServicesReq struct {
}

type ListServicesResp []string

func (c *cliHandlers) ListServicesHandler(req *ListServicesReq) ListServicesResp {
	return ListServicesResp(c.Server.services.Keys(nil))
}

func (*cliHandlers) ListServicesUsage() *handler.CommandDetails {
	return &handler.CommandDetails{
		Usage: "List services",
		Alias: "ls",
	}
}

type ListBashCommandsReq struct {
}

type ListBashCommandsResp []BashCmd

func (c *cliHandlers) ListBashCommandsHandler(req *ListServicesReq) ListBashCommandsResp {
	return ListBashCommandsResp(c.Server.BashCommands)
}

func (*cliHandlers) ListBashCommandsUsage() *handler.CommandDetails {
	return &handler.CommandDetails{
		Usage: "List bash commands",
		Alias: "lbc",
	}
}

type RunBashCommandReq struct {
	ID int
}

type RunBashCommandResp struct {
	Msg string
}

func (c *cliHandlers) RunBashCommandHandler(req *RunBashCommandReq) RunBashCommandResp {
	if req.ID < 0 || req.ID >= len(c.Server.BashCommands) {
		return RunBashCommandResp{
			Msg: "ID out of range",
		}
	}
	bc := &(c.Server.BashCommands[req.ID])
	if bc.running {
		return RunBashCommandResp{
			Msg: "Already running",
		}
	}
	bc.Run(c.Server.Commander.New())
	return RunBashCommandResp{
		Msg: "OK",
	}
}

func (*cliHandlers) RunBashCommandUsage() *handler.CommandDetails {
	return &handler.CommandDetails{
		Usage: "Run bash command",
		Alias: "rbc",
	}
}

type BashCommandOuputReq struct {
	ID int
}

type BashCommandOuputResp struct {
	Msg string
}

func (c *cliHandlers) BashCommandOuputHandler(req *RunBashCommandReq) BashCommandOuputResp {
	if req.ID < 0 || req.ID >= len(c.Server.BashCommands) {
		return BashCommandOuputResp{
			Msg: "err: ID out of range",
		}
	}
	bc := c.Server.BashCommands[req.ID]
	return BashCommandOuputResp{
		Msg: bc.buf.String(),
	}
}

func (*cliHandlers) BashCommandOuputUsage() *handler.CommandDetails {
	return &handler.CommandDetails{
		Usage: "Show output of bash command",
		Alias: "sbc",
	}
}
