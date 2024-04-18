package lusers

import (
	"sort"

	"github.com/adamcolton/luce/ds/slice"
	"github.com/adamcolton/luce/util/filter"
	"golang.org/x/crypto/bcrypt"
)

// User is a stored account: a unique ID, display Name, bcrypt
// HashedPassword and the sorted Groups it belongs to.
type User struct {
	ID             []byte `json:"-"`
	Name           string
	HashedPassword []byte
	Groups         slice.Slice[string]
}

// SetPassword replaces u's HashedPassword with the bcrypt hash of password.
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.HashedPassword = hashedPassword
	return nil
}

// CheckPassword compares password against u's HashedPassword, returning nil
// if they match (see bcrypt.CompareHashAndPassword for the error cases).
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword(u.HashedPassword, []byte(password))
}

// In reports whether u belongs to group. A nil u is never in any group.
func (u *User) In(group string) bool {
	if u == nil {
		return false
	}
	idx := u.gidx(group)
	return u.Groups.IdxCheck(idx) && u.Groups[idx] == group
}

func (u *User) gidx(group string) int {
	return u.Groups.Search(filter.GTE(group))
}

func (u *User) sortGroups() {
	sort.Slice(u.Groups, func(i, j int) bool {
		return u.Groups[i] < u.Groups[j]
	})
}

// OneRequired requires that a user be in one of the listed groups. If no groups
// are listed, the result is true.
func (u *User) OneRequired(groups []string) bool {
	if len(groups) == 0 {
		return true
	}
	if u == nil {
		return false
	}

	var group string
	fn := func(g string) bool {
		return g >= group
	}

	for _, group = range groups {
		idx := u.Groups.Search(fn)
		if u.Groups.IdxCheck(idx) && u.Groups[idx] == group {
			return true
		}
	}
	return false
}
