package lusers

import (
	"bytes"

	"github.com/adamcolton/luce/store"
)

// Group is a named set of Users, backed by its own store.NestedStore: one
// key per member, keyed by the User's ID.
type Group struct {
	Name string
	store.NestedStore
}

// We only care about the key, but we have to store some value
var hasUser = []byte{1}

// AddUser adds u as a member of g and appends g's name to u.Groups, keeping
// u.Groups sorted.
func (g *Group) AddUser(u *User) error {
	err := g.NestedStore.Put(u.ID, hasUser)
	if err != nil {
		return err
	}
	u.Groups = append(u.Groups, g.Name)
	u.sortGroups()
	return nil
}

// HasUser reports whether u is a member of g.
func (g *Group) HasUser(u *User) bool {
	return bytes.Equal(g.NestedStore.Get(u.ID).Value, hasUser)
}
