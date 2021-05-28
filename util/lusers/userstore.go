package lusers

import (
	"crypto/rand"
	"encoding/json"
	"fmt"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/store"
)

const (
	// UserIDLen is the number of random bytes UserStore.Create generates
	// for a new User's ID.
	UserIDLen = 10
	// ErrUserNotFound is returned when no User exists for the given ID or
	// name.
	ErrUserNotFound = lerr.Str("User not found")
)

var (
	byID   = []byte("server.UserStore.byID")
	byName = []byte("server.UserStore.byName")
	groups = []byte("server.GroupStore.groups")
)

// UserStore persists Users and Groups in three nested stores (by ID, by
// name, and the set of group names) built from a single store.NestedFactory.
type UserStore struct {
	byID, byName, groups store.NestedStore
}

// NewUserStore opens UserStore's three nested stores (by ID, by name and
// groups) from f.
func NewUserStore(f store.NestedFactory) (*UserStore, error) {
	us := &UserStore{}
	var err error
	us.byID, err = f.NestedStore(byID)
	if err != nil {
		return nil, err
	}
	us.byName, err = f.NestedStore(byName)
	if err != nil {
		return nil, err
	}
	us.groups, err = f.NestedStore(groups)
	if err != nil {
		return nil, err
	}
	return us, nil
}

// GetByName looks up a User by name, then loads it by ID.
func (us *UserStore) GetByName(name string) (*User, error) {
	return us.GetByID(us.byName.Get([]byte(name)).Value)
}

// Login looks up the User named name, checks password against its stored
// hash, and loads its current Groups. It returns an error if the user does
// not exist or the password does not match.
func (us *UserStore) Login(name, password string) (*User, error) {
	u, err := us.GetByID(us.byName.Get([]byte(name)).Value)
	if err != nil {
		return nil, err
	}
	err = u.CheckPassword(password)
	if err != nil {
		return nil, err
	}
	return u, us.loadGroups(u)
}

func (us *UserStore) loadGroups(u *User) error {
	us.groups.Get(nil)
	for cur := us.groups.Next(nil); cur != nil; cur = us.groups.Next(cur) {
		g, err := us.Group(string(cur))
		if err != nil {
			return err
		}
		if g.HasUser(u) {
			u.Groups = append(u.Groups, g.Name)
		}
	}
	u.sortGroups()
	return nil
}

// GetByID loads and decodes the User stored under id. It returns
// ErrUserNotFound if there is none.
func (us *UserStore) GetByID(id []byte) (*User, error) {
	b := us.byID.Get(id).Value
	if b == nil {
		return nil, ErrUserNotFound
	}
	u := &User{
		ID: id,
	}
	return u, json.Unmarshal(b, u)
}

// ErrUserAlreadyExists is returned by UserStore.Create when its name is
// already taken; the value is that name.
type ErrUserAlreadyExists string

// Error fulfills the error interface.
func (u ErrUserAlreadyExists) Error() string {
	return fmt.Sprintf("User %s already exists", string(u))
}

// Create makes a new User with the given name and password (hashed with
// bcrypt) and stores it. It returns ErrUserAlreadyExists if name is taken.
func (us *UserStore) Create(name, password string) (*User, error) {
	_, err := us.GetByName(name)
	if err == nil {
		return nil, ErrUserAlreadyExists(name)
	}
	if err != ErrUserNotFound {
		return nil, err
	}

	u := &User{
		ID:   make([]byte, UserIDLen),
		Name: name,
	}
	rand.Read(u.ID)
	err = u.SetPassword(password)
	if err != nil {
		return nil, err
	}
	err = us.Update(u)
	if err != nil {
		return nil, err
	}
	err = us.byName.Put([]byte(u.Name), u.ID)
	if err != nil {
		return nil, err
	}
	return u, nil
}

// List returns the name of every stored User.
func (us *UserStore) List() []string {
	var out []string
	for u := us.byName.Next(nil); u != nil; u = us.byName.Next(u) {
		out = append(out, string(u))
	}
	return out
}

// Update re-serializes u and stores it under its existing ID, overwriting
// the stored record. It does not touch the by-name index, so it cannot be
// used to rename a User.
func (us *UserStore) Update(u *User) error {
	b, err := json.Marshal(u)
	if err != nil {
		return err
	}
	err = us.byID.Put(u.ID, b)
	if err != nil {
		return err
	}
	return nil
}

// Groups returns the name of every Group that has been created.
func (us *UserStore) Groups() []string {
	var out []string
	for cur := us.groups.Next(nil); cur != nil; cur = us.groups.Next(cur) {
		out = append(out, string(cur))
	}
	return out
}

// Group will return a group with the provided name. If one does not already
// exist, it will be created.
func (us *UserStore) Group(name string) (*Group, error) {
	s, err := us.groups.NestedStore([]byte(name))
	if err != nil {
		return nil, err
	}
	return &Group{
		Name:        name,
		NestedStore: s,
	}, nil
}

// HasGroup will return a group with the given name only if one already exists.
func (us *UserStore) HasGroup(name string) *Group {
	r := us.groups.Get([]byte(name))
	if r.Store == nil {
		return nil
	}
	return &Group{
		Name:        name,
		NestedStore: r.Store,
	}
}
