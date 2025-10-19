package server

import (
	"fmt"
	"strings"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/cli"
	"github.com/adamcolton/luce/util/handler"
)

func (s *Server) coreCommander(ec *cli.ExitClose) cli.Commander {
	return &cliHandlers{
		Server:           s,
		ExitCloseHandler: ec.Commands(),
	}
}

// RunStdIO drives the admin CLI over the process's stdin/stdout.
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
		func(r *AdminLockUserCreationResp) {
			rnr.WriteString("Admin lock setting updated")
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

// CreateUserReq requests a new user named Name with Password.
type CreateUserReq struct {
	Name, Password string
}

// CreateUserResp reports whether CreateUserReq failed.
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

// ListUsersReq lists the server's user names.
type ListUsersReq struct{}
// ListUsersResp is the user names returned by ListUsersReq.
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

// GroupReq requests a new group named Name.
type GroupReq struct {
	Name string
}
// GroupResp reports whether GroupReq failed.
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

// ListGroupsReq lists the server's group names.
type ListGroupsReq struct{}
// ListGroupsResp is the group names returned by ListGroupsReq.
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

// UserGroupReq adds User to Group.
type UserGroupReq struct {
	User, Group string
}
// UserGroupResp reports whether UserGroupReq failed.
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

// SetPortReq changes the HTTP server's listening port.
type SetPortReq struct {
	Port string
}
// SetPortResp confirms SetPortReq; the server has already restarted.
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

// SettingsReq requests the server's current Settings.
type SettingsReq struct{}
// Settings holds the server's admin-configurable settings.
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

// AdminLockUserCreationReq sets whether new users may only be created by an
// admin.
type AdminLockUserCreationReq struct {
	AdminLockUserCreation bool
}
// AdminLockUserCreationResp confirms AdminLockUserCreationReq.
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

// RoutesReq requests the server's currently active service routes.
type RoutesReq struct {
}
// RoutesResp lists the active service routes, one per line, returned by
// RoutesReq.
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

// ListServicesReq lists the services currently registered over the service
// socket.
type ListServicesReq struct {
}

// ListServicesResp is the registered service names returned by
// ListServicesReq.
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

// ListBashCommandsReq lists the server's configured bash commands.
type ListBashCommandsReq struct {
}

// ListBashCommandsResp is the bash commands returned by ListBashCommandsReq.
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

// RunBashCommandReq runs the configured bash command at index ID.
type RunBashCommandReq struct {
	ID int
}

// RunBashCommandResp reports the result of RunBashCommandReq.
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

// BashCommandOuputReq requests the buffered output of the bash command at
// index ID.
type BashCommandOuputReq struct {
	ID int
}

// BashCommandOuputResp is the output returned by BashCommandOuputReq.
type BashCommandOuputResp struct {
	Msg string
}

func (c *cliHandlers) BashCommandOuputHandler(req *BashCommandOuputReq) BashCommandOuputResp {
	if req.ID < 0 || req.ID >= len(c.Server.BashCommands) {
		return BashCommandOuputResp{
			Msg: "err: ID out of range",
		}
	}
	bc := &(c.Server.BashCommands[req.ID])
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
