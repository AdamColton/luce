package lusers

import (
	"crypto/rand"
	"encoding/json"
	"fmt"

	"github.com/adamcolton/luce/lerr"
	"github.com/adamcolton/luce/store"
)

const (
	UserIDLen       = 10
	ErrUserNotFound = lerr.Str("User not found")
)

var (
	byID   = []byte("server.UserStore.byID")
	byName = []byte("server.UserStore.byName")
	groups = []byte("server.GroupStore.groups")
)

type UserStore struct {
	byID, byName, groups store.Store
}

func NewUserStore(f store.Factory) (*UserStore, error) {
	us := &UserStore{}
	var err error
	us.byID, err = f.Store(byID)
	if err != nil {
		return nil, err
	}
	us.byName, err = f.Store(byName)
	if err != nil {
		return nil, err
	}
	us.groups, err = f.Store(groups)
	if err != nil {
		return nil, err
	}
	return us, nil
}

func MustUserStore(f store.Factory) *UserStore {
	us, err := NewUserStore(f)
	if err != nil {
		panic(err)
	}
	return us
}

func (us *UserStore) GetByName(name string) (*User, error) {
	return us.GetByID(us.byName.Get([]byte(name)).Value)
}

func (us *UserStore) Login(name, password string) (*User, error) {
	u, err := us.GetByName(name)
	if err != nil {
		return nil, err
	}
	err = u.CheckPassword(password)
	if err != nil {
		return nil, err
	}
	return u, nil
}

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

type ErrUserAlreadyExists string

func (u ErrUserAlreadyExists) Error() string {
	return fmt.Sprintf("User %s already exists", string(u))
}

func (us *UserStore) Create(name, password string) (*User, error) {
	_, err := us.GetByName(name)
	if err == nil {
		return nil, ErrUserAlreadyExists(name)
	}
	if err != nil && err != ErrUserNotFound {
		return nil, err
	}

	u := &User{
		ID:   make([]byte, UserIDLen),
		Name: name,
	}
	rand.Read(u.ID)
	u.SetPassword(password)
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

func (us *UserStore) List() []string {
	var out []string
	for u := us.byName.Next(nil); u != nil; u = us.byName.Next(u) {
		out = append(out, string(u))
	}
	return out
}

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
	s, err := us.groups.Store([]byte(name))
	if err != nil {
		return nil, err
	}
	return &Group{
		Name:  name,
		Store: s,
	}, nil
}

// HasGroup will return a group with the given name only if one already exists.
func (us *UserStore) HasGroup(name string) *Group {
	r := us.groups.Get([]byte(name))
	if r.Store == nil {
		return nil
	}
	return &Group{
		Name:  name,
		Store: r.Store,
	}
}
