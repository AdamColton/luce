package lusess

import (
	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/util/cli"
	"github.com/adamcolton/luce/util/handler"
	"github.com/adamcolton/luce/util/lusers"
)

type StoreCmds struct {
	Store *Store
}

type CreateUserReq struct {
	Name, Password string
}

type CreateUserResp struct {
	Error error
}

func (sc *StoreCmds) UserHandler(req *CreateUserReq) *CreateUserResp {
	_, err := sc.Store.Create(req.Name, req.Password)
	return &CreateUserResp{Error: err}
}

func (sc *StoreCmds) UserUsage() *handler.CommandDetails {
	return &handler.CommandDetails{
		Usage: "Create a User",
	}
}

func CreateUserRespHandler(rnr *cli.Runner) func(resp *CreateUserResp) {
	return func(resp *CreateUserResp) {
		if resp.Error == nil {
			rnr.WriteString("user created")
		} else {
			rnr.WriteString(resp.Error.Error())
		}
	}
}

type ListUsersReq struct{}
type ListUsersResp []string

func (sc *StoreCmds) ListUsersHandler(req *ListUsersReq) ListUsersResp {
	return sc.Store.List()
}

func (sc *StoreCmds) ListUsersUsage() *handler.CommandDetails {
	return &handler.CommandDetails{
		Usage: "List User names",
	}
}

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

type GroupReq struct {
	Name string
}
type CreateGroupResp struct {
	Error error
}

func (sc *StoreCmds) GroupHandler(req *GroupReq) *CreateGroupResp {
	_, err := sc.Store.Group(req.Name)
	return &CreateGroupResp{Error: err}
}

func (sc *StoreCmds) GroupUsage() *handler.CommandDetails {
	return &handler.CommandDetails{
		Usage: "Create Group",
		Alias: "g",
	}
}

func CreateGroupRespHandler(rnr *cli.Runner) func(resp *CreateGroupResp) {
	return func(resp *CreateGroupResp) {
		if resp.Error == nil {
			rnr.WriteString("group created")
		} else {
			rnr.WriteString(resp.Error.Error())
		}
	}
}

type ListGroupsReq struct{}
type ListGroupsResp []string

func (sc *StoreCmds) ListGroupsHandler(req *ListGroupsReq) ListGroupsResp {
	return sc.Store.Groups()
}

func (sc *StoreCmds) ListGroupsUsage() *handler.CommandDetails {
	return &handler.CommandDetails{
		Usage: "List Groups",
		Alias: "lg",
	}
}

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

type UserGroupReq struct {
	User, Group string
}
type UserGroupResp struct {
	Error error
}

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
	if u == nil {
		return &UserGroupResp{
			Error: lusers.ErrUserNotFound,
		}
	}

	return &UserGroupResp{
		Error: g.AddUser(u),
	}
}

func (sc *StoreCmds) UserGroupUsage() *handler.CommandDetails {
	return &handler.CommandDetails{
		Usage: "Add User to Group",
		Alias: "ug",
	}
}

func UserGroupRespHandler(rnr *cli.Runner) func(resp *UserGroupResp) {
	return func(resp *UserGroupResp) {
		if resp.Error == nil {
			rnr.WriteString("user was added to group")
		} else {
			rnr.WriteString(resp.Error.Error())
		}
	}
}

func AllRespHandlers(rnr *cli.Runner) []any {
	return []any{
		ListUsersRespHandler(rnr),
		CreateUserRespHandler(rnr),
		CreateGroupRespHandler(rnr),
		ListGroupsRespHandler(rnr),
		UserGroupRespHandler(rnr),
	}
}
