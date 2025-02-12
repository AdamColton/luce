package lusess

import (
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/cli"
	"github.com/adamcolton/luce/util/handler"
)

// StoreCmds exposes Store's user and group management as util/handler
// commands, for a util/cli.Runner (see AllRespHandlers).
type StoreCmds struct {
	Store *Store
}

// CreateUserReq is the "user" command's request: a new User's name and
// password.
type CreateUserReq struct {
	Name, Password string
}

// CreateUserResp is the "user" command's response.
type CreateUserResp struct {
	Error error
}

// UserHandler fulfills the "user" command by creating a User.
func (sc *StoreCmds) UserHandler(req *CreateUserReq) *CreateUserResp {
	_, err := sc.Store.Create(req.Name, req.Password)
	return &CreateUserResp{Error: err}
}

// UserUsage describes the "user" command.
func (sc *StoreCmds) UserUsage() *handler.CommandDetails {
	return &handler.CommandDetails{
		Usage: "Create a User",
	}
}

// CreateUserRespHandler writes a CreateUserResp to rnr: "user created" on
// success, or the error's text.
func CreateUserRespHandler(rnr *cli.Runner) func(resp *CreateUserResp) {
	return func(resp *CreateUserResp) {
		if resp.Error == nil {
			rnr.WriteString("user created")
		} else {
			rnr.WriteString(resp.Error.Error())
		}
	}
}

// ListUsersReq is the "listUsers" command's (empty) request.
type ListUsersReq struct{}

// ListUsersResp is the "listUsers" command's response: every stored User's
// name.
type ListUsersResp []string

// ListUsersHandler fulfills the "listUsers" command.
func (sc *StoreCmds) ListUsersHandler(req *ListUsersReq) ListUsersResp {
	return sc.Store.List()
}

// ListUsersUsage describes the "listUsers" command.
func (sc *StoreCmds) ListUsersUsage() *handler.CommandDetails {
	return &handler.CommandDetails{
		Usage: "List User names",
	}
}

// ListUsersRespHandler writes each name in a ListUsersResp to rnr, one per
// line.
func ListUsersRespHandler(rnr *cli.Runner) func(users ListUsersResp) {
	return func(users ListUsersResp) {
		for i, user := range users {
			if i > 0 {
				rnr.WriteString("\n")
			}
			rnr.WriteStrings(user)
		}
	}
}

// GroupReq is the "g"/"group" command's request: the Group's name.
type GroupReq struct {
	Name string
}

// CreateGroupResp is the "g"/"group" command's response.
type CreateGroupResp struct {
	Error error
}

// GroupHandler fulfills the "g"/"group" command by creating a Group.
func (sc *StoreCmds) GroupHandler(req *GroupReq) *CreateGroupResp {
	_, err := sc.Store.Group(req.Name)
	return &CreateGroupResp{Error: err}
}

// GroupUsage describes the "g"/"group" command.
func (sc *StoreCmds) GroupUsage() *handler.CommandDetails {
	return &handler.CommandDetails{
		Usage: "Create Group",
		Alias: "g",
	}
}

// CreateGroupRespHandler writes a CreateGroupResp to rnr: "group created" on
// success, or the error's text.
func CreateGroupRespHandler(rnr *cli.Runner) func(resp *CreateGroupResp) {
	return func(resp *CreateGroupResp) {
		if resp.Error == nil {
			rnr.WriteString("group created")
		} else {
			rnr.WriteString(resp.Error.Error())
		}
	}
}

// ListGroupsReq is the "lg"/"listGroups" command's (empty) request.
type ListGroupsReq struct{}

// ListGroupsResp is the "lg"/"listGroups" command's response: every created
// Group's name.
type ListGroupsResp []string

// ListGroupsHandler fulfills the "lg"/"listGroups" command.
func (sc *StoreCmds) ListGroupsHandler(req *ListGroupsReq) ListGroupsResp {
	return sc.Store.Groups()
}

// ListGroupsUsage describes the "lg"/"listGroups" command.
func (sc *StoreCmds) ListGroupsUsage() *handler.CommandDetails {
	return &handler.CommandDetails{
		Usage: "List Groups",
		Alias: "lg",
	}
}

// ListGroupsRespHandler writes each name in a ListGroupsResp to rnr, one per
// line.
func ListGroupsRespHandler(rnr *cli.Runner) func(users ListGroupsResp) {
	return func(groups ListGroupsResp) {
		for i, group := range groups {
			if i > 0 {
				rnr.WriteString("\n")
			}
			rnr.WriteStrings(group)
		}
	}
}

// UserGroupReq is the "ug"/"userGroup" command's request: a User name and
// the Group to add it to.
type UserGroupReq struct {
	User, Group string
}

// UserGroupResp is the "ug"/"userGroup" command's response.
type UserGroupResp struct {
	Error error
}

// UserGroupHandler fulfills the "ug"/"userGroup" command by adding the
// named User to the named Group. It returns an error if either does not
// exist.
func (sc *StoreCmds) UserGroupHandler(req *UserGroupReq) *UserGroupResp {
	g := sc.Store.HasGroup(req.Group)
	if g == nil {
		return &UserGroupResp{
			Error: lerr.Str("group not found"),
		}
	}

	u, err := sc.Store.GetByName(req.User)
	if err != nil {
		return &UserGroupResp{
			Error: err,
		}
	}

	return &UserGroupResp{
		Error: g.AddUser(u),
	}
}

// UserGroupUsage describes the "ug"/"userGroup" command.
func (sc *StoreCmds) UserGroupUsage() *handler.CommandDetails {
	return &handler.CommandDetails{
		Usage: "Add User to Group",
		Alias: "ug",
	}
}

// UserGroupRespHandler writes a UserGroupResp to rnr: "user was added to
// group" on success, or the error's text.
func UserGroupRespHandler(rnr *cli.Runner) func(resp *UserGroupResp) {
	return func(resp *UserGroupResp) {
		if resp.Error == nil {
			rnr.WriteString("user was added to group")
		} else {
			rnr.WriteString(resp.Error.Error())
		}
	}
}

// AllRespHandlers returns the response handlers for all of StoreCmds'
// commands, ready to pass to a cli.Commander's Handlers.
func AllRespHandlers(rnr *cli.Runner) []any {
	return []any{
		ListUsersRespHandler(rnr),
		CreateUserRespHandler(rnr),
		CreateGroupRespHandler(rnr),
		ListGroupsRespHandler(rnr),
		UserGroupRespHandler(rnr),
	}
}
