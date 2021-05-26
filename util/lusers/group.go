package lusers

import (
	"bytes"

	"github.com/adamcolton/luce/store"
)

type Group struct {
	Name string
	store.NestedStore
}

// We only care about the key, but we have to store some value
var hasUser = []byte{1}

func (g *Group) AddUser(u *User) error {
	err := g.NestedStore.Put(u.ID, hasUser)
	if err != nil {
		return err
	}
	u.Groups = append(u.Groups, g.Name)
	u.sortGroups()
	return nil
}

func (g *Group) HasUser(u *User) bool {
	return bytes.Equal(g.NestedStore.Get(u.ID).Value, hasUser)
}
